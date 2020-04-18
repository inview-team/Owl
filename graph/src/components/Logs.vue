<template>
  <div class="container">
    <h1 v-if="currentTab == 'Alarms' ">Alarms</h1>
    <h1 v-else>Logs</h1>
    <hr>
    <ul class="nav nav-tabs" style="display: inline-flex">
      <li role="presentation" @click="getAlarms">
        <a href="#" class="btn btn-primary btn-sm" >Alarms</a>
      </li>
      <li role="presentation" @click="getLogs">
        <a href="#" class="btn btn-primary btn-sm">Logs</a>
      </li>
    </ul>
    <div class="tab-content custom-scrollbar">
      <div v-if="currentTab == 'Alarms'">
      <table class="table table-hover">
          <thead>
            <tr>
              <th scope="col">Datetime</th>
              <th scope="col">Info</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="inf in alarms" :key="inf.id">
              <td>{{ inf.time }}</td>
              <td>{{ inf.info}}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="currentTab == 'Logs'">
      <table class="table table-hover">
          <thead>
            <tr>
              <th scope="col">Datetime</th>
              <th scope="col">Info</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="inf in logs" :key="inf.id">
              <td>{{ inf.time }}</td>
              <td>{{ inf.info}}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      alarms: [],
      logs: [],
      currentTab: 'Alarms',
      api_url: process.env.RESTAPI_SVC_SERVICE_HOST,
    };
  },
  methods: {
    getAlarms() {
      this.currentTab = 'Alarms';
      
      axios.get(`http://${this.api_url}:1337/alarms`)
        .then(((res) => {
          this.alarms = res.data.alarms;
        }))
        .catch((error) => {
          console.error(error);
        });
    },
    getLogs() {
      this.currentTab = 'Logs';
      axios.get(`http://${this.api_url}:1337/logs`)
        .then(((res) => {
          this.logs = res.data.logs;
        }))
        .catch((error) => {
          console.error(error);
        });
    },
  },
  created() {
    this.getAlarms();
  },

};
</script>

<style scoped>
.custom-scrollbar {
  position: relative;
  height: 500px;
  overflow: auto;
}
.tab-content{
  display: block;
}
</style>
