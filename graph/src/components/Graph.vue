<template>
  <div class="container">
    <line-chart :chart-data="datacollection"></line-chart>
  </div>
</template>

<script>
import LineChart from './LineChart';

export default {
  components: {
    LineChart,
  },
  data() {
    return {
      datacollection: null,
    };
  },
  mounted() {
    this.fillData();
  },
  methods: {
    fillData() {
      this.datacollection = {
        labels: ['April 15th 2017', 'April 22nd 2017', 'April 24th 2017', 'April 29th 2017', 'May 1st 2017', 'May 5th 2017'],
        datasets: [
          {
            label: 'Data One',
            backgroundColor: '#f87979',
            data: [100, 500, 1000, 2000, 500, 4000],
          },
          {
            label: 'Data two',
            backgroundColor: '#008000',
            data: [4000, 500, 2000, 1000, 500, 100],
          },
        ],
      };
    },
  },
};
</script>

<style>

</style>
