<template>
  <div class="search">
    <input class="searchBox" v-model="searchTerm" placeholder="Search for presentation..." type="text"/>
    <div class="searchResults">
      <SearchResult v-for="result in results" :key="result.path + '#' + result.slide" :result="result"/>
    </div>
  </div>
</template>

<script>
import SearchResult from './SearchResult.vue'

export default {
  name: 'SearchBox',

  components: {
    SearchResult
  },

  data () {
    return {
      searchTerm: '',
      results: []
    }
  },

  watch: {
    searchTerm: function (val, oldVal) {
      let vm = this
      fetch('/api/search?q=' + val)
        .then(response => response.json())
        .then(response => {
          vm.results = response['results']
        })
    }
  }
}
</script>

<style>
.search {
  display: flex;
  flex-direction: column;
  max-width: 60rem;
  margin-left: auto;
  margin-right: auto;
}

.searchBox {
  border: 1px solid #ddd;
  box-shadow: 2px 2px 2px 0 rgba(0,0,0,0.10);
  font-size: 1.8rem;
  padding: 0.5rem;
  margin-bottom: 1rem;
  outline: none;
}

.searchResults {
  display: flex;
  flex-direction: column;
  justify-content: center;
}
</style>