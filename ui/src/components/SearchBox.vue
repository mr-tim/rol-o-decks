<template>
  <div>
    <input v-model="searchTerm" type="text"/>
    <div>
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
</style>