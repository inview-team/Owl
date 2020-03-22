<template>
  <div class="container">
    <h1>Graphs</h1>
    <p>{{metrics}}</p>
  </div>
</template>

<script>

    export default {
      data: {
        return() {
          metrics = []
        }
      },
      methods: {
        getData() {

        }
      },
      created() {
        this.getData()
      }
    }
</script>

<style scoped>

</style>
