class clsDataTable {
    template(
        id,
        name,
        description,
        headers, // Must be an array object with properties "name" for display name and "value" for the exact field name
        data, // Your array of data
        addUrl='', //Redirect url to Add page
        viewUrl='', // Example: /user/{Id}
        viewKey='',
        editUrl='', // Example: /user/{Id}/edit
        editKey='' // key used to replace the parameter on editUrl. Example: if editKey is Id, this will replace the specified {Id} to actual value for each data.
    ){
        return `
            <div x-data="{
                showAdd: ${!!addUrl ? true : false},
                showEdit: ${!!editUrl ? true : false},
                viewUrl: '${viewUrl}',
                editUrl: '${editUrl}',
                viewUrlUpdate(item){
                    return this.viewUrl.replace('{${viewKey}}', item['${viewKey}'])
                },
                editUrlUpdate(item){
                    return this.editUrl.replace('{${editKey}}', item['${editKey}'])
                },
                searchText: '',
                get filteredData(){
                    return !this.searchText ? this.${data} : this.${data}.filter(d => {
                        var found = false
                        this.${headers}.forEach(h => {
                            var colFound = d[h.value].toString().toLowerCase().includes(this.searchText.toLowerCase())
                            if (colFound) found = true; return
                        })
                        return found
                    })
                }

            }">
                <div class="mt-8 flex flex-col">
                    <div class="m-2">
                        <div class="inline-block min-w-full py-2 align-middle">
                            <div class="shadow-sm ring-1 ring-black ring-opacity-5">
                            <div class="flex flex-row m-3 pt-3">
                                <div class="sm:flex-auto">
                                    <h1 class="text-xl font-semibold text-gray-900">${name}</h1>
                                    <p class="mt-2 text-sm text-gray-700">${description}</p>
                                </div>
                                <div className="py-3">
                                    <input
                                        type="text"
                                        x-model="searchText"
                                        className="default:border-solid border-2 border-slate-300 pl-1"
                                        placeholder="Search ${name}"
                                    />
                                    <button type="button" x-show="showAdd" class="inline-flex items-center justify-center rounded-md border border-transparent bg-orange-400 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-orange-500 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 sm:w-auto" @click="window.location.href='${addUrl}'">Add</button>
                                </div>
                            </div>
                                <table id="clsDataTable_${id}" class="min-w-full border-separate" style="border-spacing: 0">
                                    <thead class="bg-gray-50">
                                        <tr>
                                            <template x-for="(h,n) in ${headers}">
                                                <th scope="col" class="sticky top-0 z-10 border-b border-gray-300 bg-gray-50 bg-opacity-75 py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter sm:pl-6 lg:pl-8" :key="n">
                                                    <span x-text="h.name"></span>
                                                </th>
                                            </template>
                                            <th scope="col" x-show="showEdit" class="sticky top-0 z-10 border-b border-gray-300 bg-gray-50 bg-opacity-75 py-3.5 pr-4 pl-3 backdrop-blur backdrop-filter sm:pr-6 lg:pr-8">
                                                <span class="sr-only">Edit</span>
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody class="bg-white">
                                        <template x-for="(i,n) in filteredData">
                                            <tr>
                                                <template x-for="(h,n) in ${headers}">
                                                <td class="border-b border-gray-200 py-4 pl-4 pr-3 text-sm text-gray-900 sm:pl-6 lg:pl-8">
                                                    <a :href="viewUrlUpdate(i)">
                                                        <span x-text="i[h.value]"></span>
                                                    </a>
                                                </td>
                                                </template>
                                                <td x-show="showEdit" class="relative whitespace-nowrap border-b border-gray-200 py-4 pr-4 pl-3 text-right text-sm font-medium sm:pr-6 lg:pr-8">
                                                    <a :href="editUrlUpdate(i)" class="inline-flex items-center justify-center rounded-md border border-transparent bg-orange-200 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-orange-500 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 sm:w-auto">Edit</a>
                                                </td>
                                            </tr>
                                        </template>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `
    }
}