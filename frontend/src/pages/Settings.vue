<template>
  <v-container>
    <v-card>
      <v-sheet class="pa-5">
        <v-form 
          ref="form"
        >
          <v-row>
            <v-col lg="2">
              Day time range:
            </v-col>
            <v-col lg="1">
              <v-text-field
                  label="from"
                  hide-details="auto"
                  v-model="settings.timeRange.min"
                  @change="updateSettings"
                  type="number"
                  :rules="rules.minTimeRange"
                  required
              ></v-text-field>
            </v-col>
            <v-col lg="1">
              <v-text-field
                  label="to"
                  hide-details="auto"
                  v-model="settings.timeRange.max"
                  @change="updateSettings"
                  type="number"
                  :rules="rules.maxTimeRange"
                  required
              ></v-text-field>
            </v-col>
          </v-row>

          Integrations:
          <v-switch
              v-model="configuration.enableIntegrations"
              inset
              label="Enable integrations with 3rd party platforms"
          ></v-switch>
          <v-simple-table v-if="configuration.enableIntegrations">
            <template v-slot:default>
              <IntegrationTable/>
            </template>
          </v-simple-table>
        </v-form>
      </v-sheet>
    </v-card>
  </v-container>
</template>

<script>
import IntegrationTable from "@/components/settings/IntegrationTable"
import SettingsService from "@/service/SettingsService"

export default {
  name: 'Settings',

  components: {
    IntegrationTable
  },
  data: () => ({
    configuration: {
      enableIntegrations: true,
    },
    timeFormat: "12H",
    settings: {
      timeRange: {},
    },
    rules: {
      minTimeRange: [
        v => !!v || 'Field is required',
        v => v > 0 || 'Should be greater than 0',
        v => v <= 24 || 'Should be less or equal to 24',
      ],
      maxTimeRange: [
        v => !!v || 'Field is required',
        v => v > 0 || 'Should be greater than 0',
        v => v <= 24 || 'Should be less or equal to 24',
      ],
    },
  }),
  methods: {
    refreshSettings() {
      SettingsService.get()
        .then(response => this.settings = response.data)
    },
    updateSettings() {
      if (this.$refs.form.validate()) {
        SettingsService.update(this.settings)
          .then(() => {
            const timeRange = this.settings.timeRange

            localStorage.minTimeRange = timeRange.min
            localStorage.maxTimeRange = timeRange.max
          })
      }
    },
  },
  created() {
    this.refreshSettings()
  },
};
</script>