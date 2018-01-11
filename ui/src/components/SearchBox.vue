<template>
  <div class="search">
    <input class="searchBox" v-model="searchTerm" type="text"/>
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
  max-width: 60rem;
  margin-left: auto;
  margin-right: auto;
}

.searchBox {
  width: 98%;
  border: 1px solid #aaa;
  box-shadow: none;
  border-radius: 0.3rem;
  font-size: 1.8rem;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  outline: none;
}

.searchResults {
  display: flex;
  flex-direction: column;
  justify-content: center;
}
</style>