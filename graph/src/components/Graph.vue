<template>
  <div class="container">
    <line-chart :chart-data="datacollection"></line-chart>
  </div>
</template>

<script>
import axios from 'axios';
import LineChart from './LineChart';

export default {
  components: {
    LineChart,
  },
  data() {
    return {
      datacollection: null,
      label: [],
      data: [],
    };
  },
  mounted() {
    this.fillData();
  },
  methods: {
    fillData() {
      axios.get('http://localhost:1337/graph')
        .then(((res) => {
          console.log(Object.keys(res.data).length);
          for (let i = 0; i < Object.keys(res.data).length; i += 1) {
            this.label.push(res.data[i].time);
            this.data.push(res.data[i].value);
          }
          this.datacollection = {
            labels: this.label,
            datasets: [
              {
                label: res.data[0].metric,
                backgroundColor: '#f87979',
                data: this.data,
              },
            ],
          };
        }))
        .catch((error) => {
          console.error(error);
        });
    },
  },
};
</script>

<style>

</style>
