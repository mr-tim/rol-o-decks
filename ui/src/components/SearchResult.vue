<template>
  <div class="searchResult">
    <div class="thumbnail">
      <img @click="open" height="200" :src="imgSrc"/>
    </div>
    <div class="matchContent" @click="open">
      <div class="matchDetails">
        <p class="presPath">{{ result.path }}</p>
        <p class="slideNum">Slide {{result.slide}}</p>
      </div>
      <div class="formattedMatch" v-html="formattedMatch"></div>
    </div>
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
.thumbnail img {
  margin: 0px;
  padding: 0px;
}

.thumbnail {
  padding: 0px;
  line-height: 0rem;
}

.searchResult {
  color: #333;
  display: flex;
  border: 1px solid #ddd;
  box-shadow: 2px 2px 2px 0 rgba(0,0,0,0.10);
  margin-bottom: 0.5rem;
  cursor: pointer;
}

.matchContent {
  display: flex;
  flex-grow: 1;
  flex-direction: column;
}

.matchDetails {
  padding-left: 0.4rem;
  margin: 0;
  border-bottom: 1px solid #ddd;
  color: #06b;
}

.formattedMatch {
  padding: 0.5rem;
}

.presPath {
  margin-top: 0.4rem;
  margin-bottom: 0.2rem;
}

.slideNum {
  margin-top: 0.2rem;
  font-size: 0.8rem;
}
</style>