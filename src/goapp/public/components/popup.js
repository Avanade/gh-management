const popup = ({
    content,
    placement = 'top', // top | right | bottom | left
    popupClass
}) => {
    return {
        content: 'CONTENT...',
        popupClass: 'absolute hidden ',
        // INITIALIZED
        async init() {
            this.content = content
            this.appendPlacementClass(placement.toLowerCase())
            this.popupClass += popupClass
        },
        hidePopup(){
            this.popupClass = this.popupClass.replace('block', 'hidden')
        },
        showPopup(){
            this.popupClass = this.popupClass.replace('hidden', 'block')
        },
        appendPlacementClass(placement) {
            switch (placement) {
                case 'top':
                    this.popupClass += "bottom-full left-1/2 transform -translate-x-1/2 bg-white p-3 rounded border mb-5 "
                break;
                case 'right':
                    this.popupClass += "top-1/2 left-full transform -translate-y-1/2 bg-white p-3 rounded border ml-3 " 
                break;
                case 'bottom':
                    this.popupClass += "top-full left-1/2 transform -translate-x-1/2 bg-white p-3 rounded border mt-2 " 
                break;
                case 'left':
                    this.popupClass += "top-1/2 right-full transform -translate-y-1/2 bg-white p-3 rounded border mr-3 " 
                break;
            }
        },
        template : `<div class="relative z-40">
                        <div x-bind:class="popupClass">
                            <div x-html="content">
                            </div>
                        </div>
                    </div>`
    }
}