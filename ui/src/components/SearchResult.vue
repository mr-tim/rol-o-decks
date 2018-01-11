<template>
  <div class="searchResult">
    <div class="thumbnail">
      <a @click="open" href="#"><img :src="imgSrc"/></a>
      <p>{{ result.path }}</p>
      <p class="slideNum">Slide {{result.slide}}</p>
    </div>
    <div class="matchContent" v-html="formattedMatch"></div>
  </div>
</template>

<script>
export default {
  name: 'SearchResult',

  props: [ 'result' ],

  methods: {
    open () {
      fetch('/api/open/' + this.result.slideId)
    }
  },

  computed: {
    fileLink () {
      return 'file://' + this.result.path
    },

    imgSrc () {
      return 'data:image/png;base64,' + this.result.thumbnail
    },

    formattedMatch () {
      let match = this.result.match
      match.end = match.start + match.length
      let output = [match.text.slice(0, match.end), '</b>', match.text.slice(match.end)].join('')
      output = [output.slice(0, match.start), '<b>', output.slice(match.start)].join('')
      return output
    }
  }
}
</script>

<style>
.slideNum {
  text-align: center;
}

.searchResult {
  display: flex;
}
</style>