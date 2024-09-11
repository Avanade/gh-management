const table = ({
    setDefault = {
      search : '',
      filter : 10,
      page : 0,
      orderby : '',
      ordertype : ''
    },
    initCallback,
    stateChangeCallback,
    callback,
    onRowClick,
    data = '',
    total = 0,
    columns = [
      // { 
      //   name : 'String', 
      //   value : 'String'|0, 
      //   render : function(),
      //   sort : {
      //     orderby : 'ColumnName' | 'ColumnIndex',
      //     ordertype : 'ASC' | 'DESC'
      //   },
      // }
    ]
  }) => {
    return { 
      columns : [],
      data : [],

      // FILTER
      search : '',
      filter : 10,
      page : 0,
      
      // ORDER BY
      orderby : '',
      ordertype : '',

      total : 0,
      showStart : 0,
      showEnd : 0,
      isLoading : false,

      async init() {
        this.filter = setDefault.filter
        this.search = setDefault.search
        this.page = setDefault.page
        this.orderby = setDefault.orderby
        this.ordertype = setDefault.ordertype

        this.columns = columns;
        await this.load();

        if (initCallback != undefined) {
          initCallback(this)
        }
      },
      async initLoad() {
        this.isLoading = true;
        this.data = [];
        this.total = 0;
        this.showStart = 0;
        this.showEnd = 0;

        this.res = await callback(this)

        this.data = this.res[data]
        this.total = this.res[total]

        this.isLoading = false;

        if (this.data == null || this.data.length == 0) return;

        this.showStart = this.data.length > 0 ? ((this.page * this.filter) + 1) : 0;
        this.showEnd = (this.page * this.filter) + this.data.length;        
      },
      async load() {
        this.isLoading = true;
        this.data = [];
        this.total = 0;
        this.showStart = 0;
        this.showEnd = 0;

        this.res = await callback(this)

        this.data = this.res[data]
        this.total = this.res[total]

        this.isLoading = false;

        if(stateChangeCallback != undefined) {
          stateChangeCallback(this)
        }

        if (this.data == null || this.data.length == 0) return;

        this.showStart = this.data.length > 0 ? ((this.page * this.filter) + 1) : 0;
        this.showEnd = (this.page * this.filter) + this.data.length;
      },
      onSetState(search, filter, page, orderby, ordertype){
        this.search = search
        this.filter = filter
        this.page = page
        this.orderby = orderby
        this.ordertype = ordertype
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
        this.filter = 10;
        this.page = 0;
        this.search = e.target.value;
        this.load();
      },
      onSortHandler(e){
        e.ordertype = e.ordertype.toLowerCase() != 'asc' ? 'ASC' : 'DESC';
        this.ordertype = e.ordertype;
        this.orderby = e.orderby;
        this.load();
      },
      onRowClickHandler(data){
        if (onRowClick == undefined) {
          console.log("ROW WAS CLICKED BUT NO EVENT WAS SET \n DATA : ", data);
          return;
        }

        onRowClick(data)
      },
      isRowClickable(){
        return onRowClick != undefined;
      },
      initRow(data){
        let html = '';
        this.columns.forEach(col => {
          for (const key in data) {
            if(key === col.value){
              html = html.concat(`<td class="whitespace-nowrap py-4 px-3 text-sm text-gray-500">${col.render != undefined ? col.render(data[key], data) : data[key]}</td>`)
            }
          }
        });
        return html;
      },
      template : `<nav class="bg-white flex items-center justify-between" aria-label="header">
                    <div class="sm:block">
                      <div>
                        <label for="filter" class="block text-sm font-medium text-gray-700">Filter</label>
                        <select @change="onChangeFilterHandler" x-model="filter" id="filter" name="filter" class="mt-1 block w-20 pl-3 pr-10 py-2 text-base text-center border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md">
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
                        <div class="mt.-1">
                          <input @keyup.enter="onSearchSubmit" type="text" name="search" id="search" class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md" x-model="search">
                        </div>
                      </div>
                    </div>
                  </nav>
                  <table class="min-w-full divide-y divide-gray-300 mt-5">
                    <thead>
                      <tr>
                        <template x-for='item in columns'>
                          <th 
                            scope="col" 
                            class="py-3.5 px-3 text-left text-sm font-semibold  text-gray-600 border-r border-gray-100 last:border-none">
                            <span class="flex group justify-between">
                                <span x-text="item.name"></span>
                                <template x-if="item.sort != undefined">
                                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 hover:text-orange-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" @click="onSortHandler(item.sort)">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
                                  </svg>
                                </template>
                            </span>
                          </th>
                        </template>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-50">
                        <template x-for='item in data'>
                          <tr x-html="initRow(item)" class="hover:bg-gray-100" :class="isRowClickable() ? 'hover:cursor-pointer' : ''" @click="onRowClickHandler(item)">
                          </tr>
                        </template>
                        <tr x-show='data == null || data.length == 0 && !isLoading' x-transition>
                          <td 
                            x-bind:colspan='columns.length'>
                            <p class="text-center my-5">NO RESULT FOUND</p>
                          </td>
                        <tr x-show='isLoading' x-transition>
                          <td x-bind:colspan='columns.length'>
                            <svg 
                              role="status" 
                              class="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-[#FF5800] m-auto my-5"
                              viewBox="0 0 100 101" 
                              fill="none" 
                              xmlns="http://www.w3.org/2000/svg">
                              <path
                                d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                                fill="currentColor" />
                              <path
                                d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                                fill="currentFill" />
                            </svg>
                          </td>
                        </tr>
                    </tbody>
                  </table>
                  <nav class="bg-white py-3 flex items-center justify-between border-t border-gray-200" aria-label="Pagination">
                    <div class="sm:block">
                      <p class="text-sm text-gray-700">
                        Showing
                        <span class="font-medium" x-text="showStart"></span>
                        to
                        <span class="font-medium" x-text="showEnd"></span>
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