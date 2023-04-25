package models

type TypPageData struct {
	Header    interface{}
	Profile   interface{}
	ProfileGH TypGitHubUser
	Content   interface{}
	HasPhoto  bool
	UserPhoto string
}

type TypGitHubUser struct {
	LoggedIn           bool
	Id                 int    `json:"id"`
	Username           string `json:"login"`
	NodeId             string `json:"node_id"`
	AvatarUrl          string `json:"avatar_url"`
	AccessToken        string
	IsValid            bool
	IsDirect           bool
	IsEnterpriseMember bool
}

type TypHeaders struct {
	Menu          []TypMenu
	ExternalLinks []TypMenu
	Page          string
}

type TypMenu struct {
	Name     string
	Url      string
	IconPath string
	External bool
}

type ExternalLinksIcon struct {
	IconName string
	IconPath string
}

type TypExternalLinks struct {
	Id                        int    `json:"id"`
	SVGName					  string `json:"svgname"`
	IconSVG                   string `json:"iconsvg"`
	Category                  string `json:"category"`	
	Created                   string `json:"created"`
	CreatedBy                 string `json:"createdBy"`
	Modified                  string `json:"modified"`
	ModifiedBy                string `json:"modifiedBy"`
}

type TypNewProjectReqBody struct {
	Id                      string `json:"id"`
	GithubId                int64  `json:"githubId"`
	Name                    string `json:"name"`
	Coowner                 string `json:"coowner"`
	Description             string `json:"description"`
	ConfirmAvaIP            bool   `json:"confirmAvaIP"`
	ConfirmSecIPScan        bool   `json:"confirmSecIPScan"`
	ConfirmNotClientProject bool   `json:"ConfirmNotClientProject"`
	TFSProjectReference     string
	Visibility              int
}

type TypeMakeProjectPublicReqBody struct {
	Id                         string `json:"id"`
	Newcontribution            string `json:"newcontribution"`
	OSSsponsor                 string `json:"osssponsor"`
	Avanadeofferingsassets     string `json:"avanadeofferingsassets"`
	Willbecommercialversion    string `json:"willbecommercialversion"`
	OSSContributionInformation string `json:"osscontributionInformation"`
}

type TypCommunity struct {
	Id                     int                   `json:"id"`
	Name                   string                `json:"name"`
	Url                    string                `json:"url"`
	Description            string                `json:"description"`
	Notes                  string                `json:"notes"`
	TradeAssocId           string                `json:"tradeAssocId"`
	IsExternal             bool                  `json:"isExternal"`
	OnBoardingInstructions string                `json:"onBoardingInstructions"`
	Created                string                `json:"created"`
	CreatedBy              string                `json:"createdBy"`
	Modified               string                `json:"modified"`
	ModifiedBy             string                `json:"modifiedBy"`
	Sponsors               []TypSponsors         `json:"sponsors"`
	Tags                   []string              `json:"tags"`
	CommunitiesExternal    []TypRelatedCommunity `json:"communitiesExternal"`
	CommunitiesInternal    []TypRelatedCommunity `json:"communitiesInternal"`
}

type TypCommunitySponsors struct {
	Id                string `json:"id"`
	CommunityId       string `json:"communityId"`
	UserPrincipalName string `json:"userprincipalname"`
	Created           string `json:"created"`
	CreatedBy         string `json:"createdBy"`
	Modified          string `json:"modified"`
	ModifiedBy        string `json:"modifiedBy"`
}
type TypProjectApprovals struct {
	Id                         int64
	ProjectId                  int64
	ProjectName                string
	ProjectCoowner             string
	ProjectDescription         string
	RequesterName              string
	RequesterGivenName         string
	RequesterSurName           string
	RequesterUserPrincipalName string
	CoownerName                string
	CoownerGivenName           string
	CoownerSurName             string
	CoownerUserPrincipalName   string
	ApprovalTypeId             int64
	ApprovalType               string
	ApproverUserPrincipalName  string
	ApprovalDescription        string
	Newcontribution            string
	OSSsponsor                 string
	Avanadeofferingsassets     string
	Willbecommercialversion    string
	OSSContributionInformation string
	RequestStatus              string
	ApproveUrl                 string
	RejectUrl                  string
	ApproveText                string
	RejectText                 string
}

type TypApprovalSystemPost struct {
	ApplicationId       string
	ApplicationModuleId string
	Email               string
	Subject             string
	Body                string
	RequesterEmail      string
}

type TypApprovalSystemPostResponse struct {
	ItemId string `json:"itemId"`
}

type TypUpdateApprovalStatusReqBody struct {
	ItemId       string `json:"itemId"`
	IsApproved   bool   `json:"isApproved"`
	Remarks      string `json:"Remarks"`
	ResponseDate string `json:"responseDate"`
}

type TypSponsors struct {
	DisplayName string `json:"displayName"`
	Mail        string `json:"mail"`
}
type TypRelatedCommunities struct {
	Name       string `json:"Name"`
	Url        string `json:"Url"`
	IsExternal bool   `json:"IsExternal"`
}

type TypCommunitySponsorsList struct {
	Name      string `json:"Name"`
	GivenName string `json:"GivenName"`
	SurName   string `json:"SurName"`
	Email     string `json:"Email"`
}

type TypCommunityOnBoarding struct {
	Id                     int64                      `json:"Id"`
	Name                   string                     `json:"Name"`
	Url                    string                     `json:"Url"`
	OnBoardingInstructions string                     `json:"OnBoardingInstructions"`
	Sponsors               []TypCommunitySponsorsList `json:"Sponsors"`
	Communities            []TypRelatedCommunities    `json:"Communities"`
}

type TypCommunityApprovals struct {
	Id                         int64
	CommunityId                int64
	CommunityName              string
	CommunityUrl               string
	CommunityDescription       string
	CommunityNotes             string
	CommunityTradeAssocId      string
	CommunityIsExternal        bool
	RequesterName              string
	RequesterGivenName         string
	RequesterSurName           string
	RequesterUserPrincipalName string
	ApproverUserPrincipalName  string
	ApprovalDescription        string
}

type TypCategory struct {
	Id               int                   `json:"id"`
	Name             string                `json:"name"`
	Created          string                `json:"created"`
	CreatedBy        string                `json:"createdBy"`
	Modified         string                `json:"modified"`
	ModifiedBy       string                `json:"modifiedBy"`
	CategoryArticles []TypCategoryArticles `json:"categoryArticles"`
}

type TypCategoryArticles struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Url          string `json:"Url"`
	Body         string `json:"Body"`
	CategoryId   int    `json:"CategoryId"`
	CategoryName string `json:"CategoryName"`
	Created      string `json:"created"`
	CreatedBy    string `json:"createdBy"`
	Modified     string `json:"modified"`
	ModifiedBy   string `json:"modifiedBy"`
}
type TypCommunityApprovers struct {
	Id                        int    `json:"id"`
	ApproverUserPrincipalName string `json:"name"`
	Disabled                  bool   `json:"disabled"`
	Created                   string `json:"created"`
	CreatedBy                 string `json:"createdBy"`
	Modified                  string `json:"modified"`
	ModifiedBy                string `json:"modifiedBy"`
}

type TypRelatedCommunity struct {
	ParentCommunityId  int `json:"ParentCommunityId"`
	RelatedCommunityId int `json:"RelatedCommunityId"`
}

type TypUpdateApprovalReAssign struct {
	Id                  string `json:"id"`
	ApproverEmail       string `json:"ApproverEmail"`
	Username            string `json:"Username"`
	ApplicationId       string `json:"ApplicationId"`
	ApplicationModuleId string `json:"ApplicationModuleId"`
	ItemId              string `json:"itemId"`
	ApproveText         string `json:"ApproveText"`
	RejectText          string `json:"RejectText"`
}
