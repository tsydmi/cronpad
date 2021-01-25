<template>
  <v-container>
    <v-card
        tile
        min-height="400px"
        class="pa-10"
    >
      <v-col class="d-flex justify-end">
        <v-tooltip right>
          <template v-slot:activator="{ on, attrs }">
            <v-btn
                icon
                class="primary--text"
                @click="createProjectDialog = true"
                v-bind="attrs"
                v-on="on"
            >
              <v-icon>mdi-plus-circle</v-icon>
            </v-btn>
          </template>
          <span>ADD PROJECT</span>
        </v-tooltip>
      </v-col>

      <span v-if="!projects || projects.length === 0" class="inactive--text mt-10 justify-center d-flex">
        No projects found
      </span>

      <v-data-table
          :headers="projectsHeaders"
          :items="projects"
          hide-default-footer
      >
        <template v-slot:item.color="{ item }">
          <v-icon v-bind:color="item.color">mdi-format-color-highlight</v-icon>
        </template>

        <template v-slot:item.actions="{ item }">
          <div class="d-flex justify-end">
            <v-btn
                icon
                @click="selectProject(item)"
            >
              <v-icon>mdi-pencil-outline</v-icon>
            </v-btn>
            <v-btn
                icon
                @click="openDeleteProjectDialog(item)"
            >
              <v-icon>mdi-trash-can-outline</v-icon>
            </v-btn>
          </div>
        </template>

        <template v-slot:item.users="{ item }">
          <span> {{ getAllUserNames(item) }}</span>
        </template>

      </v-data-table>

      <create-project-dialog v-model="createProjectDialog"
                             :users="users"
                             @refreshProjects="refreshProjects"/>

      <edit-project-dialog v-model="editProjectDialog"
                           :users="users"
                           :project="projectToEdit"
                           @refreshProjects="refreshProjects"/>

      <delete-dialog v-model="deleteProjectDialog"
                     @delete="deleteSelectedProject"
                     :item-name="projectToDelete.name"/>
    </v-card>
  </v-container>
</template>

<script>
import CreateProjectDialog from "@/components/projects/CreateProjectDialog"
import EditProjectDialog from "@/components/projects/EditProjectDialog"
import ProjectService from "@/service/ProjectService"
import UserService from "@/service/UserService"
import DeleteDialog from "@/components/DeleteDialog"
import cloneDeep from 'clone-deep'

export default {
  name: "Projects",
  components: {
    CreateProjectDialog,
    EditProjectDialog,
    DeleteDialog,
  },

  data: () => ({
    createProjectDialog: false,

    editProjectDialog: false,
    projectToEdit: {},

    deleteProjectDialog: false,
    projectToDelete: {},

    projects: [],
    projectsHeaders: [
      {text: 'NAME', value: 'name'},
      {text: 'DESCRIPTION', value: 'description'},
      {text: 'USERS', value: 'users', sortable: false},
      {text: '', value: 'actions', sortable: false},
    ],
    users: [],
    search: {},
  }),
  methods: {
    openDeleteProjectDialog(project) {
      this.projectToDelete = project
      this.deleteProjectDialog = true
    },
    deleteSelectedProject() {
      ProjectService.delete(this.projectToDelete.id)
          .then(() => this.refreshProjects())
    },
    clearFilters() {
      this.search = {}
      this.refreshProjects()
    },
    selectProject(project) {
      this.projectToEdit = cloneDeep(project)
      this.editProjectDialog = true
    },
    getAllUserNames(project) {
      if (!project.users) {
        return ''
      }

      return project.users
          .map(u => this.getFullNameById(u))
          .join(', ')
    },
    getFullNameById(id) {
      let user = this.users.find(u => u.id === id);
      if (user) {
        return `${user.firstName} ${user.lastName}`
      }
      return "error error"
    },
    refreshProjects() {
      ProjectService.search(this.search)
          .then(response => this.projects = response.data)
    },
    refreshUsers() {
      UserService.findAll()
          .then((response) => this.users = response.data)
    },
  },
  created() {
    this.refreshProjects()
    this.refreshUsers()
  },
}
</script>