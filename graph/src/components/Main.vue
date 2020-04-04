<template>
<div class="container">
  <graph></graph>
  <div style="display: flex">
  <settings></settings>
  <logs></logs>
  </div>
</div>
</template>

<script>
import Logs from './Logs';
import Settings from './Settings';
import Graph from './Graph';

export default {
  name: 'Main',
  components: {
    logs: Logs,
    settings: Settings,
    graph: Graph,
  },
};
</script>

<style scoped>

</style>
