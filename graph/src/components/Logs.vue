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
              <td>{{ inf.datetime }}</td>
              <td>{{ inf.Info}}</td>
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
              <td>{{ inf.datetime }}</td>
              <td>{{ inf.Info}}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>
</template>

<script>
    export default {
       data() {
         return {
           alarms: [{
             "id": 0,
             "datetime": "22.03.2020",
             "Info": "Alarm"
           }],
           logs: [{
             "id": 0,
             "datetime": "21.03.2020",
             "Info": "Logs"
           }
           ],
           chooser: ['Alarms','Log'],
           currentTab: "Alarms",
         }
       },
      methods: {
        getAlarms() {

        },
        getLogs(){

        }
      }

    }
</script>

<style scoped>

</style>
