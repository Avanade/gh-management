const deeplink  = () => {
    return {
        url: {
            // hash: '',
            // host: '',
            // hostname: '',
            // href: '',
            // origin: '',
            // pathname: '',
            // protocol: '',
            // search: '',
            // searchParams: {},
        },
        // INITIALIZED
        async init() {
            this.url = new URL(window.location.href)
        },
        onSetPopsetEventListener(data, callback) {
            window.addEventListener('popstate', function(e){
                console.log('TRIGGER POPSTATE EVENT')
                data.url = new URL(e.target.location.href)
                callback(data)
            });
        },
        // EVENT HANDLER
        onSetParams(params) {
            console.log('SET PARAMS')
            // params { name : 'name', value : 'value' }
            params.forEach(param => {
                if(param.value == '' || param.value == undefined) {
                    this.url.searchParams.delete(param.name)
                }
                else {
                    this.url.searchParams.set(param.name, param.value)
                }
            });
            this.onPushState()
        },
        onPushState() {
            console.log('PUSH STATE')
            let urlPath = this.url.origin

            if (this.url.pathname != '/') {
                urlPath = `${urlPath}${this.url.pathname}`
            }

            window.history.pushState({}, '', `${urlPath}${this.url.search}`)
        }
    }
}