package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pusher/oauth2_proxy/pkg/apis/sessions"
	"github.com/pusher/oauth2_proxy/pkg/http_cache"
)

// EpitechProvider represents an Epitech Intranet based Identity Provider
type EpitechProvider struct {
	*AzureProvider
	AuthToken string

	// GroupValidator is a function that determines if the passed email is in
	// the configured Epitech group.
	GroupValidator func(string) bool

	client *http.Client
}

type epitechGroupMember struct {
	Type     string `json:"type"`
	Login    string `json:"login"`
	Slug     string `json:"slug"`
	Location string `json:"location"`
	Title    string `json:"title"`
	Close    bool   `json:"close"`
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

		//Create a custom client so we can make use of our RoundTripper
		//If you make use of http.Get(), the default http client located at http.DefaultClient is used instead
		//Since we have special needs, we have to make use of our own http.RoundTripper implementation
		client: &http.Client{
			Transport: http_cache.NewCacheTransport(http.DefaultTransport, 60),
			Timeout:   time.Second * 5,
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
	// email, err := p.GetEmailAddress(s)
	// if err != nil {
	// 	return nil, err
	// }
	// rawUser, err := p.getEpitechUser(email, p.AuthToken)
	// if err != nil {
	// 	return nil, err
	// }
	// s.User = rawUser

	return
}

// ValidateGroup validates that the provided email exists in the configured Epitech
// group(s).
func (p *EpitechProvider) ValidateGroup(email string) bool {
	return p.GroupValidator(email)
}

func (p *EpitechProvider) verifyGroupMembership(groups []string, authToken string, email string) bool {
	fmt.Println("Checking Epitech group membership for " + email)

	for _, group := range groups {
		members, err := p.getEpitechGroupMembers(group, authToken)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		for _, member := range members {
			if member.Slug == email || member.Login == email {
				return true
			}
		}
	}

	return false
}

func (p *EpitechProvider) getEpitechGroupMembers(groupName string, authToken string) ([]epitechGroupMember, error) {
	path := fmt.Sprintf("https://intra.epitech.eu/%s/group/%s/member?format=json", authToken, groupName)

	var req *http.Request
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []epitechGroupMember{}, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return []epitechGroupMember{}, err
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return []epitechGroupMember{}, err
	}

	var members []epitechGroupMember
	err = json.Unmarshal([]byte(body), &members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (p *EpitechProvider) getEpitechUser(email string, authToken string) (string, error) {
	path := fmt.Sprintf("https://intra.epitech.eu/%s/user/%s/?format=json", authToken, email)

	var req *http.Request
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Do(req)
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
