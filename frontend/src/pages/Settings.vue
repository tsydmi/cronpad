<template>
  <v-container>
    <v-card>
      <v-sheet class="pa-5">
        <v-form>
          <v-radio-group
              v-model="timeFormat"
              row
              label="Time format:"
              disabled
          >
            <v-radio label="12 hours" value="12H"></v-radio>
            <v-radio label="24 hours" value="24H"></v-radio>
          </v-radio-group>

          <v-row>
            <v-col lg="2">
              Day time range:
            </v-col>
            <v-col lg="1">
              <v-text-field
                  label="from"
                  hide-details="auto"
                  value="7 AM"
              ></v-text-field>
            </v-col>
            <v-col lg="1">
              <v-text-field
                  label="to"
                  hide-details="auto"
                  value="11 PM"
              ></v-text-field>
            </v-col>
          </v-row>

          <v-divider inset class="pb-5"></v-divider>

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
import IntegrationTable from "@/components/settings/IntegrationTable";

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
  }),
};
</script>