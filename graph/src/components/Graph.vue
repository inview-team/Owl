<template>
  <div class="container">
    <b-form-select :options="metrics"
                   style="width: 30%"
                   v-model="selected"
                   v-on:change="fillData"
    ></b-form-select>
    <line-chart :chart-data="datacollection"></line-chart>
  </div>
</template>

<script>
import axios from 'axios';
import items from '@/data/data.json';
import LineChart from './LineChart';


export default {
  components: {
    LineChart,
  },
  data() {
    return {
      datacollection: null,
      times: [],
      data: [],
      metrics: [],
      selected: '',
      api_url: items.url,
    };
  },
  mounted() {
    this.fillData();
  },
  methods: {
    fillData() {
      this.metrics = [];
      this.data = [];
      this.times = [];
      axios.get(`http://${this.api_url}:1337/graph`)
        .then(((res) => {
          console.log(Object.keys(res.data).length);
          for (let i = 0; i < Object.keys(res.data).length; i += 1) {
            if (this.metrics.indexOf(res.data[i].metric) === -1) {
              this.metrics.push(res.data[i].metric);
              this.data.push([]);
              this.times.push([]);
            }
            this.times[this.metrics.indexOf(res.data[i].metric)].push(res.data[i].time);
            this.data[this.metrics.indexOf(res.data[i].metric)].push(res.data[i].value);
          }
          if (this.selected === '') {
            this.selected = this.metrics[0];
          }
          this.datacollection = {
            labels: this.times[this.metrics.indexOf(this.selected)],
            datasets: [
              {
                label: this.metrics[this.metrics.indexOf(this.selected)],
                backgroundColor: '#f87979',
                data: this.data[this.metrics.indexOf(this.selected)],
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
