const table = ({
    callback,
    data = '',
    total = 0,
    columns = [
      // { Name : 'String', Value : 'String'|0 }
    ]
  }) => {
    return { 
      columns : [],
      data : [],
      search : '',
      filter : 10,
      page : 0,
      total : 0,
      isLoading : false,
      async init() {
        this.columns = columns;
        await this.load();
      },
      async load() {
        this.isLoading = true;
        this.data = [];
        this.total = 0;

        this.res = await callback({
          filter : this.filter,
          page : this.page,
          search : this.search
        })

        this.data = this.res[data]
        this.total = this.res[total]

        this.isLoading = false;
      },
      nextPageEnabled(){
        return this.page < Math.ceil(this.total/this.filter) - 1
      },
      onNextPageHandler(){
        if (!this.nextPageEnabled()) return;

        this.page = this.page + 1

        this.load()
      },
      previousPageEnabled(){
        return this.page > 0
      },
      onPreviousPageHandler(){
        if (!this.previousPageEnabled) return;

        this.page = this.page - 1

        this.load();
      },
      onChangeFilterHandler(e){
        this.filter = parseInt(e.target.value);
        this.load()
      },
      onSearchSubmit(e){
        this.search = e.target.value;
        this.load();
      },
      initRow(data){
        let html = '';
        this.columns.forEach(col => {
          for (const key in data) {
            if(key === col.value){
              html = html.concat(`<td class="whitespace-nowrap py-4 px-3 text-sm text-gray-500">${data[key]}</td>`) 
            }
          }
        });
        return html;
      },
      template : `<nav class="bg-white flex items-center justify-between" aria-label="header">
                    <div class="sm:block">
                      <div>
                        <label for="filter" class="block text-sm font-medium text-gray-700">Filter</label>
                        <select @change="onChangeFilterHandler" id="filter" name="filter" class="mt-1 block w-20 pl-3 pr-10 py-2 text-base text-center border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md">
                          <option>5</option>
                          <option>10</option>
                          <option>20</option>
                          <option>50</option>
                          <option>100</option>
                        </select>
                      </div>
                    </div>
                    <div class="flex justify-between sm:justify-end">
                      <div class="sm:col-span-3">
                        <label for="search" class="block text-sm font-medium text-gray-700">Search</label>
                        <div class="mt-1">
                          <input @keyup.enter="onSearchSubmit" type="text" name="search" id="search" class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md" x-model="search">
                        </div>
                      </div>
                    </div>
                  </nav>
                  <table class="min-w-full divide-y divide-gray-300">
                    <thead>
                      <tr>
                        <!-- HEADER HERE -->
                        <template x-for='item in columns'>
                          <th scope="col" class="py-3.5 px-3 text-left text-sm font-semibold text-gray-900" x-text="item.name"></th>
                        </template>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-200">
                        <template x-for='item in data'>
                          <tr x-html="initRow(item)">
                          </tr>
                        </template>
                        <tr x-show='isLoading' x-transition>
                          <td x-bind:colspan='columns.length'>
                            <div class="flex justify-center items-center">
                              <span>Loading...</span>
                            </div>
                          </td>
                        </tr>
                    </tbody>
                  </table>
                  <nav class="bg-white py-3 flex items-center justify-between border-t border-gray-200" aria-label="Pagination">
                    <div class="sm:block">
                      <p class="text-sm text-gray-700">
                        Showing
                        <span class="font-medium" x-text="data.length > 0 ? ((page * filter) + 1) : 0"></span>
                        to
                        <span class="font-medium" x-text="((page * filter) + data.length)"></span>
                        of
                        <span class="font-medium" x-text="total"></span>
                        results
                      </p>
                    </div>
                    <div class="flex justify-between sm:justify-end">
                      <button x-bind:disabled="!previousPageEnabled()" x-on:click="onPreviousPageHandler" href="#" class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-200"> Previous </button>
                      <button x-bind:disabled="!nextPageEnabled()" x-on:click="onNextPageHandler" href="#" class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-200"> Next </button>
                    </div>
                  </nav>`
    }
  }