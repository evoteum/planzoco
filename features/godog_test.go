package features_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/evoteum/planzoco/databases"
	"github.com/evoteum/planzoco/routes"

	"github.com/cucumber/godog"
)

var serverURL string

// TestFeatures wires godog up to `go test`, running every *.feature file in this
// directory against a real instance of the app over HTTP, backed by a real Postgres
// reachable via DATABASE_URL (same variable the app itself uses).
func TestFeatures(t *testing.T) {
	if err := databases.InitDB(); err != nil {
		t.Fatalf("failed to init database (is DATABASE_URL set to a reachable Postgres?): %v", err)
	}

	// Capture the features directory (containing the .feature files) before chdir'ing away,
	// so godog still finds them regardless of what directory routes.SetupRoutes needs below.
	featuresDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// routes.SetupRoutes loads templates/static assets relative to the app module root
	// (go/planzoco), but `go test` runs with this package's directory as the working directory.
	if err := os.Chdir("../go/planzoco"); err != nil {
		t.Fatalf("failed to chdir to app module root: %v", err)
	}

	server := httptest.NewServer(routes.SetupRoutes())
	defer server.Close()
	serverURL = server.URL

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{featuresDir},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type ctxKey struct{}

// state holds everything a scenario accumulates as it drives the app: the shared
// "browser" used for actions that don't care about voter identity, one browser per
// named voter (each with its own cookie jar, simulating a separate person), and the
// IDs/text discovered along the way so later steps can refer back to "the question" etc.
type state struct {
	client *http.Client
	voters map[string]*http.Client

	eventID       string
	questionIDs   map[string]string // question text -> id
	optionIDs     map[string]string // option text -> id
	currentQText  string

	lastStatus int
	lastBody   string
}

func newState() *state {
	return &state{
		client:      newClient(),
		voters:      map[string]*http.Client{},
		questionIDs: map[string]string{},
		optionIDs:   map[string]string{},
	}
}

func newClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{Jar: jar}
}

func (s *state) voter(name string) *http.Client {
	if c, ok := s.voters[name]; ok {
		return c
	}
	c := newClient()
	s.voters[name] = c
	return c
}

func stateFrom(ctx context.Context) *state {
	return ctx.Value(ctxKey{}).(*state)
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, ctxKey{}, newState()), nil
	})

	sc.Step(`^an event named "([^"]*)"$`, iCreateAnEventNamed)
	sc.Step(`^I create an event named "([^"]*)"$`, iCreateAnEventNamed)
	sc.Step(`^I rename the event to "([^"]*)"$`, iRenameTheEventTo)
	sc.Step(`^I delete the event$`, iDeleteTheEvent)
	sc.Step(`^I create an event with a (\d+) character name$`, iCreateAnEventWithNCharName)

	sc.Step(`^the event has a question "([^"]*)"$`, iAddAQuestion)
	sc.Step(`^I add a question "([^"]*)" to the event$`, iAddAQuestion)
	sc.Step(`^I add a question "([^"]*)" to the event as "([^"]*)"$`, iAddAQuestionAs)
	sc.Step(`^I add a question with (\d+) characters to the event$`, iAddAQuestionWithNChars)
	sc.Step(`^I change the question to "([^"]*)"$`, iChangeTheQuestionTo)
	sc.Step(`^I delete the question$`, iDeleteTheQuestion)

	sc.Step(`^the question has an option "([^"]*)"$`, iSuggestAnOption)
	sc.Step(`^I suggest the option "([^"]*)" for the question$`, iSuggestAnOption)
	sc.Step(`^I suggest an option with (\d+) characters for the question$`, iSuggestAnOptionWithNChars)
	sc.Step(`^I change the option "([^"]*)" to "([^"]*)"$`, iChangeTheOptionTo)
	sc.Step(`^I delete the option "([^"]*)"$`, iDeleteTheOption)

	sc.Step(`^I have voted for "([^"]*)" as "([^"]*)"$`, iVoteForAs)
	sc.Step(`^I vote for "([^"]*)" as "([^"]*)"$`, iVoteForAs)
	sc.Step(`^I vote for "([^"]*)" without giving a name$`, iVoteWithoutAName)

	sc.Step(`^I request the health endpoint$`, iRequestTheHealthEndpoint)

	sc.Step(`^the event should exist$`, theEventShouldExist)
	sc.Step(`^the event should no longer exist$`, theEventShouldNoLongerExist)
	sc.Step(`^the question should no longer exist$`, theQuestionShouldNoLongerExist)
	sc.Step(`^the option should no longer exist$`, theOptionShouldNoLongerExist)
	sc.Step(`^the option "([^"]*)" should no longer exist$`, theNamedOptionShouldNoLongerExist)
	sc.Step(`^the event should be rejected$`, theRequestShouldBeRejected)
	sc.Step(`^the question should be rejected$`, theRequestShouldBeRejected)
	sc.Step(`^the option should be rejected$`, theRequestShouldBeRejected)

	sc.Step(`^the event page should show "([^"]*)"$`, theEventPageShouldShow)
	sc.Step(`^the event page should list the question "([^"]*)"$`, theEventPageShouldShow)
	sc.Step(`^the question page should list the option "([^"]*)"$`, theQuestionPageShouldShow)
	sc.Step(`^the question should now read "([^"]*)"$`, theQuestionPageShouldShow)
	sc.Step(`^the event page should show "No votes yet" for the question$`, theEventPageShouldShowNoVotesYetForTheQuestion)
	sc.Step(`^the event page should show "([^"]*)" as the leading answer$`, theEventPageShouldShowAsTheLeadingAnswer)
	sc.Step(`^the question should invite adding an option$`, theQuestionShouldInviteAddingAnOption)
	sc.Step(`^the question should invite voting$`, theQuestionShouldInviteVoting)
	sc.Step(`^"([^"]*)" should have (\d+) vote$`, optionShouldHaveVotes)
	sc.Step(`^"([^"]*)" should have (\d+) votes$`, optionShouldHaveVotes)
	sc.Step(`^I should be asked who I am$`, iShouldBeAskedWhoIAm)
	sc.Step(`^I should land on the question's own page$`, iShouldLandOnTheQuestionsOwnPage)
	sc.Step(`^the question page should show that "([^"]*)" asked the question$`, theQuestionPageShouldShowAskedBy)
	sc.Step(`^the event page should show that "([^"]*)" asked the question$`, theEventPageShouldShowAskedBy)
	sc.Step(`^the question page should show "([^"]*)" voted for "([^"]*)"$`, theQuestionPageShouldShowVotedFor)

	sc.Step(`^the response status should be (\d+)$`, theResponseStatusShouldBe)
	sc.Step(`^the response body should be "([^"]*)"$`, theResponseBodyShouldBe)
}

// --- HTTP helpers ---

func post(client *http.Client, path string, values url.Values) (*http.Response, string, error) {
	resp, err := client.PostForm(serverURL+path, values)
	if err != nil {
		return nil, "", fmt.Errorf("POST %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return resp, string(body), nil
}

func get(client *http.Client, path string) (*http.Response, string, error) {
	resp, err := client.Get(serverURL + path)
	if err != nil {
		return nil, "", fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return resp, string(body), nil
}

// --- scraping helpers: the app is server-rendered HTML with no JSON API, so
// scenarios recover IDs the same way a real browser would - by following links
// and reading the rendered page. ---

var optionRowRe = regexp.MustCompile(`(?s)<p class="option-text">([^<]*)</p>\s*<span class="votes">Votes: (\d+)(?: \(([^)]*)\))?</span>\s*<form action="/options/([^/]+)/vote"[^>]*>.*?<a href="([^"]*)" class="edit-link">Edit</a>\s*<form action="([^"]*)" method="POST"[^>]*>\s*<button type="submit" class="delete-link">Delete</button>`)

type optionRow struct {
	ID           string
	Votes        int
	Voters       string
	EditHref     string
	DeleteAction string
}

func parseOptions(html string) map[string]optionRow {
	result := map[string]optionRow{}
	for _, m := range optionRowRe.FindAllStringSubmatch(html, -1) {
		votes, _ := strconv.Atoi(m[2])
		result[m[1]] = optionRow{ID: m[4], Votes: votes, Voters: m[3], EditHref: m[5], DeleteAction: m[6]}
	}
	return result
}

// extractHref finds an <a class="edit-link"> whose text is exactly linkText, and
// returns its href. Used to prove a UI entry point actually exists, rather than
// assuming a route is reachable just because the handler for it exists.
func extractHref(html, linkText string) (string, bool) {
	re := regexp.MustCompile(`<a href="([^"]*)" class="edit-link">` + regexp.QuoteMeta(linkText) + `</a>`)
	m := re.FindStringSubmatch(html)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// extractDeleteFormAction finds a delete <form>/<button class="delete-link"> pair
// whose button text is exactly buttonText, and returns the form's action.
func extractDeleteFormAction(html, buttonText string) (string, bool) {
	re := regexp.MustCompile(`(?s)<form action="([^"]*)" method="POST"[^>]*>\s*<button type="submit" class="delete-link">` + regexp.QuoteMeta(buttonText) + `</button>`)
	m := re.FindStringSubmatch(html)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// extractSaveFormAction finds the "Save Changes" form shared by all three edit
// pages (edit_event.html, edit_question.html, edit_option.html) and returns its action.
func extractSaveFormAction(html string) (string, bool) {
	re := regexp.MustCompile(`(?s)<form action="([^"]*)" method="POST">.*?<button type="submit">Save Changes</button>`)
	m := re.FindStringSubmatch(html)
	if m == nil {
		return "", false
	}
	return m[1], true
}

var questionRowRe = regexp.MustCompile(`(?s)<div class="question-text">([^<]*)<span class="asked-by">asked by ([^<]*)</span></div>\s*<div class="answer-text">(.*?)</div>\s*<div class="row-actions">\s*<a href="/questions/([^"]+)" class="vote-link">([^<]*)</a>`)

func parseQuestions(html string) map[string]struct {
	ID       string
	AskedBy  string
	Answer   string
	LinkText string
} {
	result := map[string]struct {
		ID       string
		AskedBy  string
		Answer   string
		LinkText string
	}{}
	for _, m := range questionRowRe.FindAllStringSubmatch(html, -1) {
		result[m[1]] = struct {
			ID       string
			AskedBy  string
			Answer   string
			LinkText string
		}{ID: m[4], AskedBy: m[2], Answer: strings.TrimSpace(m[3]), LinkText: m[5]}
	}
	return result
}

func repeatChar(n int) string {
	return strings.Repeat("a", n)
}

// --- Given/When steps ---

func iCreateAnEventNamed(ctx context.Context, name string) (context.Context, error) {
	st := stateFrom(ctx)
	resp, body, err := post(st.client, "/events", url.Values{"name": {name}})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body

	if resp.StatusCode == http.StatusOK && strings.HasPrefix(resp.Request.URL.Path, "/events/") {
		st.eventID = strings.TrimPrefix(resp.Request.URL.Path, "/events/")
	}
	return ctx, nil
}

func iCreateAnEventWithNCharName(ctx context.Context, n int) (context.Context, error) {
	return iCreateAnEventNamed(ctx, repeatChar(n))
}

func iRenameTheEventTo(ctx context.Context, name string) (context.Context, error) {
	st := stateFrom(ctx)
	if st.eventID == "" {
		return ctx, fmt.Errorf("no event to rename")
	}

	_, eventBody, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return ctx, err
	}
	editHref, ok := extractHref(eventBody, "Edit Event")
	if !ok {
		return ctx, fmt.Errorf(`no "Edit Event" link found on the event page`)
	}

	_, editBody, err := get(st.client, editHref)
	if err != nil {
		return ctx, err
	}
	saveAction, ok := extractSaveFormAction(editBody)
	if !ok {
		return ctx, fmt.Errorf("no save form found on the edit-event page")
	}

	resp, body, err := post(st.client, saveAction, url.Values{"name": {name}})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

func iDeleteTheEvent(ctx context.Context) (context.Context, error) {
	st := stateFrom(ctx)
	if st.eventID == "" {
		return ctx, fmt.Errorf("no event to delete")
	}

	_, eventBody, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return ctx, err
	}
	deleteAction, ok := extractDeleteFormAction(eventBody, "Delete Event")
	if !ok {
		return ctx, fmt.Errorf(`no "Delete Event" button found on the event page`)
	}

	resp, body, err := post(st.client, deleteAction, url.Values{})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

// addQuestion posts a new question as client, transparently handling the
// whoami name gate (using fallbackName the first time this client needs one)
// the same way a real browser would: gated, name itself, then retry.
func addQuestion(st *state, client *http.Client, eventID, text, fallbackName string) (*http.Response, string, error) {
	resp, body, err := post(client, "/events/"+eventID+"/questions", url.Values{"text": {text}})
	if err != nil {
		return nil, "", err
	}

	if resp.Request.URL.Path == "/whoami" {
		next := resp.Request.URL.Query().Get("next")
		if _, _, err := post(client, "/whoami", url.Values{"name": {fallbackName}, "next": {next}}); err != nil {
			return nil, "", err
		}
		resp, body, err = post(client, "/events/"+eventID+"/questions", url.Values{"text": {text}})
		if err != nil {
			return nil, "", err
		}
	}

	// Adding a question redirects straight to the new question's own page,
	// so its ID comes from where we landed, not from scraping the event page.
	if resp.StatusCode == http.StatusOK && strings.HasPrefix(resp.Request.URL.Path, "/questions/") {
		st.questionIDs[text] = strings.TrimPrefix(resp.Request.URL.Path, "/questions/")
		st.currentQText = text
	}
	return resp, body, nil
}

func iAddAQuestion(ctx context.Context, text string) (context.Context, error) {
	st := stateFrom(ctx)
	if st.eventID == "" {
		return ctx, fmt.Errorf("no event to add a question to")
	}
	resp, body, err := addQuestion(st, st.client, st.eventID, text, "Tester")
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

func iAddAQuestionAs(ctx context.Context, text, name string) (context.Context, error) {
	st := stateFrom(ctx)
	if st.eventID == "" {
		return ctx, fmt.Errorf("no event to add a question to")
	}
	resp, body, err := addQuestion(st, st.voter(name), st.eventID, text, name)
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

func iAddAQuestionWithNChars(ctx context.Context, n int) (context.Context, error) {
	return iAddAQuestion(ctx, repeatChar(n))
}

func iChangeTheQuestionTo(ctx context.Context, newText string) (context.Context, error) {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return ctx, err
	}

	_, qBody, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return ctx, err
	}
	editHref, ok := extractHref(qBody, "Edit Question")
	if !ok {
		return ctx, fmt.Errorf(`no "Edit Question" link found on the question page`)
	}

	_, editBody, err := get(st.client, editHref)
	if err != nil {
		return ctx, err
	}
	saveAction, ok := extractSaveFormAction(editBody)
	if !ok {
		return ctx, fmt.Errorf("no save form found on the edit-question page")
	}

	resp, body, err := post(st.client, saveAction, url.Values{"text": {newText}})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body

	delete(st.questionIDs, st.currentQText)
	st.questionIDs[newText] = qID
	st.currentQText = newText
	return ctx, nil
}

func iDeleteTheQuestion(ctx context.Context) (context.Context, error) {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return ctx, err
	}

	_, qBody, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return ctx, err
	}
	deleteAction, ok := extractDeleteFormAction(qBody, "Delete Question")
	if !ok {
		return ctx, fmt.Errorf(`no "Delete Question" button found on the question page`)
	}

	resp, body, err := post(st.client, deleteAction, url.Values{})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

func (s *state) currentQuestionID() (string, error) {
	id, ok := s.questionIDs[s.currentQText]
	if !ok {
		return "", fmt.Errorf("no current question")
	}
	return id, nil
}

func iSuggestAnOption(ctx context.Context, text string) (context.Context, error) {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return ctx, err
	}
	resp, body, err := post(st.client, "/questions/"+qID+"/options", url.Values{"text": {text}})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body

	if o, ok := parseOptions(body)[text]; ok {
		st.optionIDs[text] = o.ID
	}
	return ctx, nil
}

func iSuggestAnOptionWithNChars(ctx context.Context, n int) (context.Context, error) {
	return iSuggestAnOption(ctx, repeatChar(n))
}

func iChangeTheOptionTo(ctx context.Context, oldText, newText string) (context.Context, error) {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return ctx, err
	}

	_, qBody, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return ctx, err
	}
	opt, ok := parseOptions(qBody)[oldText]
	if !ok {
		return ctx, fmt.Errorf("option %q not found on question page", oldText)
	}
	if opt.EditHref == "" {
		return ctx, fmt.Errorf("no edit link found for option %q", oldText)
	}

	_, editBody, err := get(st.client, opt.EditHref)
	if err != nil {
		return ctx, err
	}
	saveAction, ok := extractSaveFormAction(editBody)
	if !ok {
		return ctx, fmt.Errorf("no save form found on the edit-option page")
	}

	resp, body, err := post(st.client, saveAction, url.Values{"text": {newText}})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body

	delete(st.optionIDs, oldText)
	st.optionIDs[newText] = opt.ID
	return ctx, nil
}

func iDeleteTheOption(ctx context.Context, text string) (context.Context, error) {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return ctx, err
	}

	_, qBody, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return ctx, err
	}
	opt, ok := parseOptions(qBody)[text]
	if !ok {
		return ctx, fmt.Errorf("option %q not found on question page", text)
	}
	if opt.DeleteAction == "" {
		return ctx, fmt.Errorf("no delete button found for option %q", text)
	}

	resp, body, err := post(st.client, opt.DeleteAction, url.Values{})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

func iVoteForAs(ctx context.Context, optionText, voterName string) (context.Context, error) {
	st := stateFrom(ctx)
	optID, ok := st.optionIDs[optionText]
	if !ok {
		return ctx, fmt.Errorf("unknown option %q", optionText)
	}
	client := st.voter(voterName)

	resp, body, err := post(client, "/options/"+optID+"/vote", url.Values{})
	if err != nil {
		return ctx, err
	}

	if resp.Request.URL.Path == "/whoami" {
		next := resp.Request.URL.Query().Get("next")
		if _, _, err := post(client, "/whoami", url.Values{"name": {voterName}, "next": {next}}); err != nil {
			return ctx, err
		}
		resp, body, err = post(client, "/options/"+optID+"/vote", url.Values{})
		if err != nil {
			return ctx, err
		}
	}

	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

func iVoteWithoutAName(ctx context.Context, optionText string) (context.Context, error) {
	st := stateFrom(ctx)
	optID, ok := st.optionIDs[optionText]
	if !ok {
		return ctx, fmt.Errorf("unknown option %q", optionText)
	}

	resp, body, err := post(newClient(), "/options/"+optID+"/vote", url.Values{})
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

func iRequestTheHealthEndpoint(ctx context.Context) (context.Context, error) {
	st := stateFrom(ctx)
	resp, body, err := get(st.client, "/health")
	if err != nil {
		return ctx, err
	}
	st.lastStatus = resp.StatusCode
	st.lastBody = body
	return ctx, nil
}

// --- Then steps ---

func theEventShouldExist(ctx context.Context) error {
	st := stateFrom(ctx)
	if st.eventID == "" {
		return fmt.Errorf("no event was created (last status %d)", st.lastStatus)
	}
	resp, _, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected event to exist, got status %d", resp.StatusCode)
	}
	return nil
}

func theEventShouldNoLongerExist(ctx context.Context) error {
	st := stateFrom(ctx)
	resp, _, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected event to be gone, got status %d", resp.StatusCode)
	}
	return nil
}

func theQuestionShouldNoLongerExist(ctx context.Context) error {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return err
	}
	resp, _, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected question to be gone, got status %d", resp.StatusCode)
	}
	return nil
}

func theOptionShouldNoLongerExist(ctx context.Context) error {
	st := stateFrom(ctx)
	var optID string
	for _, id := range st.optionIDs {
		optID = id
	}
	if optID == "" {
		return fmt.Errorf("no option to check")
	}
	return optionShouldBeGone(st, optID)
}

func theNamedOptionShouldNoLongerExist(ctx context.Context, text string) error {
	st := stateFrom(ctx)
	optID, ok := st.optionIDs[text]
	if !ok {
		return fmt.Errorf("no option named %q to check", text)
	}
	return optionShouldBeGone(st, optID)
}

func optionShouldBeGone(st *state, optID string) error {
	resp, _, err := get(st.client, "/options/"+optID+"/edit")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected option to be gone, got status %d", resp.StatusCode)
	}
	return nil
}

func theRequestShouldBeRejected(ctx context.Context) error {
	st := stateFrom(ctx)
	if st.lastStatus != http.StatusBadRequest {
		return fmt.Errorf("expected status 400, got %d", st.lastStatus)
	}
	return nil
}

func theEventPageShouldShow(ctx context.Context, text string) error {
	st := stateFrom(ctx)
	if st.eventID == "" {
		return fmt.Errorf("no event to check")
	}
	_, body, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	if !strings.Contains(body, text) {
		return fmt.Errorf("expected event page to contain %q, it did not", text)
	}
	return nil
}

func theQuestionPageShouldShow(ctx context.Context, text string) error {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return err
	}
	_, body, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return err
	}
	if !strings.Contains(body, text) {
		return fmt.Errorf("expected question page to contain %q, it did not", text)
	}
	return nil
}

func theEventPageShouldShowNoVotesYetForTheQuestion(ctx context.Context) error {
	st := stateFrom(ctx)
	_, body, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	q, ok := parseQuestions(body)[st.currentQText]
	if !ok {
		return fmt.Errorf("question %q not found on event page", st.currentQText)
	}
	if q.Answer != "No votes yet" {
		return fmt.Errorf("expected %q, got %q", "No votes yet", q.Answer)
	}
	return nil
}

func theEventPageShouldShowAsTheLeadingAnswer(ctx context.Context, optionText string) error {
	st := stateFrom(ctx)
	_, body, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	q, ok := parseQuestions(body)[st.currentQText]
	if !ok {
		return fmt.Errorf("question %q not found on event page", st.currentQText)
	}
	if !strings.Contains(q.Answer, optionText) {
		return fmt.Errorf("expected leading answer to contain %q, got %q", optionText, q.Answer)
	}
	return nil
}

func theQuestionShouldInviteAddingAnOption(ctx context.Context) error {
	st := stateFrom(ctx)
	_, body, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	q, ok := parseQuestions(body)[st.currentQText]
	if !ok {
		return fmt.Errorf("question %q not found on event page", st.currentQText)
	}
	if strings.Contains(strings.ToLower(q.LinkText), "vote") {
		return fmt.Errorf("expected the question to invite adding an option, link read %q", q.LinkText)
	}
	return nil
}

func theQuestionShouldInviteVoting(ctx context.Context) error {
	st := stateFrom(ctx)
	_, body, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	q, ok := parseQuestions(body)[st.currentQText]
	if !ok {
		return fmt.Errorf("question %q not found on event page", st.currentQText)
	}
	if !strings.Contains(strings.ToLower(q.LinkText), "vote") {
		return fmt.Errorf("expected the question to invite voting, link read %q", q.LinkText)
	}
	return nil
}

func optionShouldHaveVotes(ctx context.Context, optionText string, want int) error {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return err
	}
	_, body, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return err
	}
	opt, ok := parseOptions(body)[optionText]
	if !ok {
		return fmt.Errorf("option %q not found on question page", optionText)
	}
	if opt.Votes != want {
		return fmt.Errorf("expected %q to have %d votes, got %d", optionText, want, opt.Votes)
	}
	return nil
}

func theQuestionPageShouldShowAskedBy(ctx context.Context, name string) error {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return err
	}
	_, body, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return err
	}
	if !strings.Contains(body, "asked by "+name) {
		return fmt.Errorf("expected question page to show %q asked the question, body did not contain it", name)
	}
	return nil
}

func theEventPageShouldShowAskedBy(ctx context.Context, name string) error {
	st := stateFrom(ctx)
	_, body, err := get(st.client, "/events/"+st.eventID)
	if err != nil {
		return err
	}
	q, ok := parseQuestions(body)[st.currentQText]
	if !ok {
		return fmt.Errorf("question %q not found on event page", st.currentQText)
	}
	if q.AskedBy != name {
		return fmt.Errorf("expected event page to show %q asked the question, got %q", name, q.AskedBy)
	}
	return nil
}

func theQuestionPageShouldShowVotedFor(ctx context.Context, voterName, optionText string) error {
	st := stateFrom(ctx)
	qID, err := st.currentQuestionID()
	if err != nil {
		return err
	}
	_, body, err := get(st.client, "/questions/"+qID)
	if err != nil {
		return err
	}
	opt, ok := parseOptions(body)[optionText]
	if !ok {
		return fmt.Errorf("option %q not found on question page", optionText)
	}
	if !strings.Contains(opt.Voters, voterName) {
		return fmt.Errorf("expected %q to show %q as a voter, got voters %q", optionText, voterName, opt.Voters)
	}
	return nil
}

func iShouldBeAskedWhoIAm(ctx context.Context) error {
	st := stateFrom(ctx)
	if !strings.Contains(st.lastBody, "What does this group know you as?") {
		return fmt.Errorf("expected to be shown the name prompt, was not")
	}
	return nil
}

func iShouldLandOnTheQuestionsOwnPage(ctx context.Context) error {
	st := stateFrom(ctx)
	if !strings.Contains(st.lastBody, "Add your suggestions and vote on them!") {
		return fmt.Errorf("expected to land on the question page, did not")
	}
	return nil
}

func theResponseStatusShouldBe(ctx context.Context, want int) error {
	st := stateFrom(ctx)
	if st.lastStatus != want {
		return fmt.Errorf("expected status %d, got %d", want, st.lastStatus)
	}
	return nil
}

func theResponseBodyShouldBe(ctx context.Context, want string) error {
	st := stateFrom(ctx)
	if strings.TrimSpace(st.lastBody) != want {
		return fmt.Errorf("expected body %q, got %q", want, st.lastBody)
	}
	return nil
}
