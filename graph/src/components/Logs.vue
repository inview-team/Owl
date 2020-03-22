<template>
  <div class="container">
    <h1>Logs</h1>
    <ul class="nav nav-tabs">
      <li role="presentation" @click="currentTab='Alarms'">
        <a href="#">Alarms</a>
      </li>
      <li role="presentation" @click="currentTab='Logs'">
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
            axios.get(`http://localhost:3000/alarm`)
                .then( (res => {
                  this.alarms = res.data;
                }))
                .catch( (error) => {
                  console.error(error);
                })
        },
        getLogs(){
            axios.get(`http://localhost:3000/logs`)
                .then( (res => {
                  this.logs = res.data;
                }))
                .catch( (error) => {
                  console.error(error);
                })
        }
      },
      created() {
         this.getAlarms()
         this.getLogs()
      }

    }
</script>

<style scoped>

</style>
