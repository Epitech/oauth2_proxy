package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pusher/oauth2_proxy/pkg/apis/sessions"
)

// EpitechProvider represents an Epitech Intranet based Identity Provider
type EpitechProvider struct {
	*AzureProvider
	AuthToken string

	// GroupValidator is a function that determines if the passed email is in
	// the configured Epitech group.
	GroupValidator func(string) bool
}

type epitechUserInfoGroup struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type epitechUserInfo struct {
	Groups []epitechUserInfoGroup `json:"groups"`
}

// NewEpitechProvider initiates a new EpitechProvider
func NewEpitechProvider(p *ProviderData) *EpitechProvider {
	provider := &EpitechProvider{
		AzureProvider: NewAzureProvider(p),

		// Set a default GroupValidator to just always return valid (true), it will
		// be overwritten if we configured a Epitech group restriction.
		GroupValidator: func(email string) bool {
			return true
		},
	}

	provider.ProviderName = "Epitech"

	return provider
}

// Configure defaults the EpitechProvider configuration options
func (p *EpitechProvider) Configure(tenant string, authToken string, groups []string) {
	p.AzureProvider.Configure(tenant)
	p.AuthToken = authToken

	if len(groups) > 0 {
		p.GroupValidator = func(email string) bool {
			return p.verifyGroupMembership(groups, authToken, email)
		}
	}
}

// GetEmailAddress returns the Account email address
func (p *EpitechProvider) GetEmailAddress(s *sessions.SessionState) (string, error) {
	return p.AzureProvider.GetEmailAddress(s)
}

// GetUserName returns the Account email address
func (p *EpitechProvider) GetUserName(s *sessions.SessionState) (string, error) {
	return p.AzureProvider.GetEmailAddress(s)
}

func (p *EpitechProvider) Redeem(redirectURL, code string) (s *sessions.SessionState, err error) {
	s, err = p.AzureProvider.Redeem(redirectURL, code)
	if err != nil {
		return
	}

	// commented for now (probably too much data for X-Forwarded-User header)
	email, err := p.GetEmailAddress(s)
	if err != nil {
		return nil, err
	}
	rawUser, err := p.getEpitechUser(email, p.AuthToken)
	if err != nil {
		return nil, err
	}
	s.User = rawUser

	return
}

// ValidateGroup validates that the provided email exists in the configured Epitech
// group(s).
func (p *EpitechProvider) ValidateGroup(email string) bool {
	return p.GroupValidator(email)
}

func (p *EpitechProvider) verifyGroupMembership(groups []string, authToken string, email string) bool {
	fmt.Println("Checking Epitech group membership for " + email)

	rawUser, err := p.getEpitechUser(email, authToken)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	user, err := p.unserializeEpitechUser(rawUser)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, group := range groups {
		for _, userGroup := range user.Groups {
			if userGroup.Name == group {
				return true
			}
		}
	}

	return false
}

func (p *EpitechProvider) getEpitechUser(email string, authToken string) (string, error) {
	path := fmt.Sprintf("https://intra.epitech.eu/%s/user/%s/?format=json", authToken, email)

	var req *http.Request
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", err
	}
	return string(body), nil
}

func (p *EpitechProvider) unserializeEpitechUser(body string) (*epitechUserInfo, error) {
	var user epitechUserInfo

	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
