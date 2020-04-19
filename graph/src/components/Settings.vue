<template>
  <div class="container">
    <h1>Settings</h1>
    <hr>
    <button class="btn btn-success btn-sm" v-b-modal.settings-modal>New Settings</button>
    <div>
      <table class="table table-hover">
          <thead>
            <tr>
              <th scope="col">Metric</th>
              <th scope="col">From</th>
              <th scope="col">To</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="inf in info" :key="inf.id">
              <td>{{ inf.metric }}</td>
              <td>{{ inf.from}}</td>
              <td>{{ inf.to}}</td>
            </tr>
          </tbody>
        </table>
    </div>
       <b-modal ref="editSettingsModal"
           id="settings-modal"
           title="Update Settings"
           hide-footer>
        <b-form @submit="onSubmit" @reset="onReset">
          <b-form-group id="form-metric-input"
                  label="Choose metric"
                  label-for="form-metric-input">
            <b-form-select
                  :options="variants"
                  v-model="editSettingsForm.metric"
            ></b-form-select>
          </b-form-group>

          <b-form-group id="form-from-group"
                    label="From:"
                    label-for="form-from-input">
            <b-form-input id="form-from-input"
                        type="text"
                        v-model="editSettingsForm.from"
                        required
                        placeholder="Enter from">
            </b-form-input>
          </b-form-group>

          <b-form-group id="form-to-group"
                    label="To:"
                    label-for="form-to-input">
            <b-form-input id="form-to-input"
                        type="text"
                        v-model="editSettingsForm.to"
                        required
                        placeholder="Enter to">
            </b-form-input>
          </b-form-group>
          <b-button type="submit" variant="primary">Submit</b-button>
          <b-button type="reset" variant="danger">Reset</b-button>
        </b-form>
       </b-modal>
  </div>
</template>

<script>
import axios from 'axios';
import items from '@/data/data.json';

export default {
  data() {
    return {
      test: [],
      info: [],
      variants: ['Pressure', 'Humidity', 'Room Temperature', 'Working area Temperature', 'pH', 'Weight', 'Fluid flow', 'CO2'],
      editSettingsForm: {
        metric: '',
        from: null,
        to: null,
      },
      api_url: items.url,
    };
  },
  methods: {
    getInfo() {
      axios.get(`http://${this.api_url}:1337/settings`)
        .then(((res) => {
          this.info = res.data.settings;
        }))
        .catch((error) => {
          console.error(error);
        });
    },
    initForm() {
      this.editSettingsForm.metric = '';
      this.editSettingsForm.from = null;
      this.editSettingsForm.to = null;
    },
    updateSettings(payload, id) {
      axios.put(`http://${this.api_url}:1337/settings_update/${id}`, payload)
        .then((responce) => {
          this.getInfo();
          console.log(responce);
        })
        .catch((error) => {
          console.log(error);
          this.getInfo();
        });
      const today = new Date();
      const date = `${today.getFullYear()}-${today.getMonth() + 1}-${today.getDate()}`;
      const time = `${today.getHours()}:${today.getMinutes()}:${today.getSeconds()}`;
      const dateTime = `${date} ${time}`;
      axios.post(`http://${this.api_url}:1337/logs`, {
        time: dateTime,
        info: `Update ${this.editSettingsForm.metric}`,
      });
    },
    onSubmit(evt) {
      evt.preventDefault();
      this.$refs.editSettingsModal.hide();
      console.log(this.editSettingsForm.metric);
      const metrId = this.info.find(inf => inf.metric === this.editSettingsForm.metric).id;
      const payload = {
        id: metrId,
        metric: this.editSettingsForm.metric,
        from: this.editSettingsForm.from,
        to: this.editSettingsForm.to,
      };
      this.updateSettings(payload, metrId);
      this.initForm();
    },
    onReset(evt) {
      evt.preventDefault();
      this.$refs.editSettingsModal.hide();
      this.initForm();
    },
  },
  created() {
    this.getInfo();
  },
};
</script>
