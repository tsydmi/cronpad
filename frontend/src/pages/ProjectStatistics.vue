<template>
  <v-container>
    <v-card
        class="pa-4"
        min-height="600px"
    >
      <v-row>
        <v-col cols="3">
          <v-list flat v-if="projects && projects.length > 0">
            <v-list-item-group
                color="primary"
                @change="selectProject"
            >
              <v-list-item
                  v-for="(project, i) in projects"
                  :key="i"
              >
                <v-list-item-content
                >
                  <v-list-item-title> {{ project.name }}</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
          <span class="ma-5" v-else>
            No projects found
          </span>
        </v-col>

        <v-col cols="9" class="pl-10 pr-10">
          <div class="row" v-if="selectedProject !== null">
            <v-col cols="12">
              <p class="text-lg-h3"> {{ selectedProject.name }} </p>
            </v-col>

            <v-col cols="12">
              <div v-if="users.length > 0">
                Users:
                <span class="pl-5">{{ getAllUserNames(users) }}</span>
              </div>
              <span v-else class="pl-5"> No users assigned to selected project</span>
            </v-col>

            <v-col cols="6">
              <progress-chart
                  v-model="progressChartData"
                  element-id="progress-chart"/>
              <div class="d-flex justify-center">
                <span v-if="selectedProject">{{ getProjectTime(selectedProject) }}</span>
              </div>
            </v-col>
            <v-col cols="6">
              <pie-chart
                  v-model="tagChartData"
                  element-id="tag-chart"
                  :sum="hoursSum"
                  title="Tags"/>
            </v-col>

            <v-col cols="12">
              <v-data-table
                  :headers="tagTableHeaders"
                  :items="tagTableData"
                  hide-default-footer
              >
                <template v-slot:item.parent="{ item }">
                  <span v-if="item.parent" :style="{color: item.parent.color}">
                    {{ item.parent.name }}
                  </span>
                </template>
              </v-data-table>
            </v-col>
          </div>

          <div class="ml-10 pt-10 inactive--text" v-else>
            <span>
              No project selected
            </span>
          </div>
        </v-col>
      </v-row>
    </v-card>
  </v-container>
</template>

<script>
import ReportService from "@/service/ReportService"
import ProjectService from "@/service/ProjectService"
import PieChart from "@/components/PieChart"
import ProgressChart from "@/components/ProgressChart"
import dayjs from "dayjs"

export default {
  name: "ProjectStatistics",

  components: {
    PieChart,
    ProgressChart,
  },

  data: () => ({
    hasAdminRole: false,
    selectedProject: null,
    projects: [],

    users: [],
    hoursSum: null,

    tagTableHeaders: [
      {text: 'Tag Name', value: 'name'},
      {text: 'Parent', value: 'parent'},
      {text: 'Hours', value: 'hours'},
      {text: 'Percent', value: 'percent', sortable: false}
    ],
    tagTableData: [],

    tagChartData: null,
    progressChartData: null,
  }),
  methods: {
    getAllUserNames(users) {
      if (!users) {
        return ''
      }

      return users
          .map(u => `${u.firstName} ${u.lastName}`)
          .join(', ')
    },
    selectProject(index) {
      if (index === undefined || index === null) {
        this.selectedProject = null
      } else {
        this.selectedProject = this.projects[index]
        this.getStatistics(this.selectedProject)
      }
    },
    getStatistics(project) {
      ProjectService.getUsers(project.id)
          .then(response => this.users = response.data)

      ReportService.getProjectReport(project.id)
          .then(response => {
            this.tagChartData = response.data.tagChart
            this.tagTableData = response.data.tagSummaryTable
            this.progressChartData = {
              daysPassed: response.data.daysPassed,
              daysAhead: response.data.daysAhead,
            }
            this.hoursSum = response.data.hoursSum
          })
    },
    getProjects() {
      ProjectService.findAll(this.hasAdminRole)
          .then(response => this.projects = response.data)
    },
    getProjectTime(project) {
      if (project.start && project.end) {
        return `${dayjs(project.start).format('YYYY-MM-DD')} - ${dayjs(project.end).format('YYYY-MM-DD')}`
      } else {
        if (project.start) {
          return `${dayjs(project.start).format('YYYY-MM-DD')} - now`
        }
      }

      return ''
    },
  },
  created() {
    this.hasAdminRole = this.$keycloak.hasRealmRole('admin')
    this.getProjects()
  },
}
</script>