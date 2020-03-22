<template>
  <div class="container">
    <h1 v-if="currentTab == 'Alarms' ">Alarms</h1>
    <h1 v-else>Logs</h1>
    <ul class="nav nav-tabs">
      <li role="presentation" @click="getAlarms">
        <a href="#">Alarms</a>
      </li>
      <li role="presentation" @click="getLogs">
        <a href="#">Logs</a>
      </li>
    </ul>
    <div class="tab-content">
      <div v-if="currentTab == 'Alarms'">
      <table class="table table-hover">
          <thead>
            <tr>
              <th scope="col">Datetime</th>
              <th scope="col">Info</th>
              <th></th>
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
              <th></th>
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
           currentTab: "Alarms",
         }
       },
      methods: {
        getAlarms() {
          this.currentTab = "Alarms";

            axios.get(`http://localhost:5000/alarms`)
                .then( (res => {
                  this.alarms = res.data.alarms;
                }))
                .catch( (error) => {
                  console.error(error);
                })
        },
        getLogs(){
          this.currentTab = "Logs"
            axios.get(`http://localhost:5000/logs`)
                .then( (res => {
                  this.logs = res.data.logs;
                }))
                .catch( (error) => {
                  console.error(error);
                })
        }
      },
      created() {
         this.getAlarms();
      }

    }
</script>

<style scoped>

</style>
