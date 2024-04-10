const message = (
    isInnersourceMember,
    isOpensourceOrgMember,
    innersourceOrg,
    openSourceOrg,
    isInvalidToken
) => {
    return {
        isInvalidToken: false,
        isInnersourceMember: false,
        isOpensourceOrgMember: false,
        innersourceOrg: "",
        openSourceOrg: "",
        init() {
            console.log(isInvalidToken, isInnersourceMember, isOpensourceOrgMember, innersourceOrg, openSourceOrg)
            this.isInvalidToken = isInvalidToken;
            this.isInnersourceMember = isInnersourceMember;
            this.isOpensourceOrgMember = isOpensourceOrgMember;
            this.innersourceOrg = innersourceOrg;
            this.openSourceOrg = openSourceOrg;
        },
        isDisplayMessage() {
            return !this.isInnersourceMember || !this.isOpensourceOrgMember
        },
        displayMessage() {
            const innersourceOrgLink = this.getOrgInvitationLink(this.innersourceOrg)
            const opensourceOrgLink = this.getOrgInvitationLink(this.openSourceOrg)

            if (this.isInvalidToken) {
                return `We currently can't check your membership on ${innersourceOrgLink} and ${opensourceOrgLink}. Please try again later.`
            }

            if (!this.isInnersourceMember && !this.isOpensourceOrgMember) {
                return `You have a pending invitation to ${innersourceOrgLink} and ${opensourceOrgLink}. You need to join these organizations before you can request for a repository.`
            } else if (!this.isInnersourceMember) {
                return `You have a pending invitation to ${innersourceOrgLink}. You need to join these organizations before you can request for a repository.`
            } else if (!this.isOpensourceMember) {
                return `You have a pending invitation to ${opensourceOrgLink}. You need to join these organizations before you can request for a repository.`
            }
        },
        getOrgInvitationLink(orgName) {
            return `<a href="https://github.com/orgs/${orgName}/invitation" target="_blank"><b>${orgName}</b></a>`
        },
        template: `<template x-if="isDisplayMessage">
                        <div class="bg-red-50 p-4 mb-3">
                            <div class="flex">
                                <div class="flex-shrink-0">
                                    <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                                    </svg>
                                </div>
                                <div class="ml-3 flex-1 md:flex md:justify-between">
                                    <p class="text-sm text-red-700">
                                        <span x-html="displayMessage"></span>
                                    </p>
                                </div>
                            </div>
                        </div>
                    </template>`
    }
}