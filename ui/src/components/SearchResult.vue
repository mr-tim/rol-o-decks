<template>
<div>
  <a @click="open" href="#"><img :src="imgSrc"/></a>
  <p>{{ result.path }}</p>
  <p>Slide {{result.slide}}</p>
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
    }
  }
}
</script>

<style>

</style>