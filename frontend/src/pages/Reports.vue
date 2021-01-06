<template>
  <v-container>
    <v-card
        class="pa-4"
        min-height="600px"
    >
      <v-row>
        <v-col cols="3">
          <user-section @change="selectUser" v-model="selectedUser"/>
        </v-col>
        <v-col cols="9" class="pl-10 pr-10">
          <v-form :disabled="selectedUser === null" class="row">
            <v-spacer/>
            <v-col cols="4">
              <date-range-picker v-model="dateRange" @change="refreshReports"/>
            </v-col>
          </v-form>

          <div class="row" v-if="selectedUser !== null">
            <v-col cols="12">
              <p class="text-lg-h3"> {{ selectedUser.firstName }} {{ selectedUser.lastName }} </p>
              <span>Summary: {{ hoursSum }} hours </span>
            </v-col>
            <v-col cols="6">
              <pie-chart element-id="tag-chart" v-model="tagChartData" title="Tags"/>
            </v-col>
            <v-col cols="6">
              <pie-chart element-id="project-chart" v-model="projectChartData" title="Projects"/>
            </v-col>
            <v-col cols="12">
              <date-range-bar-chart element-id="data-range-chart" v-model="dateRangeChartData"/>
            </v-col>
            <v-col cols="12">
              <v-data-table
                  :headers="eventsTableHeaders"
                  :items="eventsTableData"
                  hide-default-footer
              ></v-data-table>
            </v-col>
          </div>

          <div class="ml-10 inactive--text" v-else>
            <span>
              No user selected
            </span>
          </div>
        </v-col>
      </v-row>
    </v-card>
  </v-container>
</template>

<script>
import ReportService from "@/service/ReportService"
import UserSection from "@/components/reports/UserSection"
import PieChart from "@/components/PieChart"
import DateRangeBarChart from "@/components/reports/DateRangeBarChart"
import DateRangePicker from "@/components/DateRangePicker"
import dayjs from 'dayjs'

export default {
  name: "Reports",

  components: {
    UserSection,
    PieChart,
    DateRangeBarChart,
    DateRangePicker,
  },

  data: () => ({
    dateRange: [dayjs().date(1).format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')],
    selectedUser: null,

    hoursSum: 0,
    groupedEvents: [],

    tagChartData: null,
    projectChartData: null,
    dateRangeChartData: null,

    eventsTableHeaders: [
      {text: 'Event Name', value: 'name'},
      {text: 'Hours', value: 'hours'},
      {text: 'Percent', value: 'percent', sortable: false}
    ],
    eventsTableData: [],
  }),
  methods: {
    selectUser(user) {
      this.selectedUser = user

      this.refreshReports()
    },
    refreshReports() {
      let search = {}

      if (this.selectedUser) {
        search.userID = this.selectedUser.id

        if (this.dateRange) {
          if (this.dateRange.length > 0) {
            search.from = dayjs(this.dateRange[0], 'YYYY-MM-DD')
                .hour(0).minute(0).second(0)
                .toISOString()

            if (this.dateRange.length === 2) {
              search.to = dayjs(this.dateRange[1], 'YYYY-MM-DD')
                  .hour(23).minute(59).second(59)
                  .toISOString()
            }
          }
        }

        this.getReports(search)
      }
    },
    getReports(query) {
      ReportService.findAll(query)
          .then(response => {
            this.tagChartData = response.data.tagChart
            this.projectChartData = response.data.projectChart
            this.dateRangeChartData = response.data.dateRangeChartDataSets;

            this.eventsTableData = response.data.eventSummaryTable
            this.hoursSum = response.data.hoursSum
          });
    },
  },
}
</script>