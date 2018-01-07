<template>
<div>
  <a target="_blank" :href="fileLink">
    <img :src="imgSrc"/>
  </a>
</div>
</template>

<script>
export default {
  name: 'SearchResult',

  props: [ 'result' ],

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