<template>
  <v-container>
    <v-card
        tile
        min-height="400px"
        class="pa-10"
    >
      <v-form>
        <v-row class="mb-3">
          <v-col cols="2">
            <v-text-field
                label="NAME"
                v-model="search.name"
                @keydown.enter="refreshProjects"
                @change="refreshProjects"
                dense
            ></v-text-field>
          </v-col>

          <v-col cols="3">
            <v-text-field
                label="DESCRIPTION"
                v-model="search.description"
                @keydown.enter="refreshProjects"
                @change="refreshProjects"
                dense
            ></v-text-field>
          </v-col>

          <v-col cols="5">
            <user-select v-model="search.users" :users="users"
                         label="USERS"
                         @change="refreshProjects">
              <template v-slot:selection="{fullName, index}">
                <span class="mr-2" v-if="index < 2">
                  {{ fullName }}
                </span>
                <span v-if="index === 2" class="grey--text caption">
                  (+{{ search.users.length - 2 }} others)
                </span>
              </template>
            </user-select>
          </v-col>

          <v-col cols="1">
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                    icon
                    @click="clearFilters"
                    v-bind="attrs"
                    v-on="on"
                    :disabled="!search.name && !search.description && (!search.users || search.users.length === 0)"
                >
                  <v-icon>far fa-times-circle</v-icon>
                </v-btn>
              </template>
              <span>Clear filters</span>
            </v-tooltip>
          </v-col>

          <v-col cols="1" class="d-flex justify-end">
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
        </v-row>
      </v-form>

      <span v-if="!projects || projects.length === 0" class="inactive--text mt-10 justify-center d-flex">
        No projects found
      </span>

      <v-row
          v-for="project in projects"
          :key="project.id">
        <v-col cols="2">{{ project.name }}</v-col>
        <v-col cols="3">{{ project.description }}</v-col>
        <v-col cols="5">
          <span>{{ getAllUserNames(project) }}</span>
        </v-col>
        <v-col cols="2" class="d-flex justify-end">
          <v-btn
              icon
              @click="selectProject(project)"
          >
            <v-icon>mdi-pencil-outline</v-icon>
          </v-btn>
          <v-btn
              icon
              @click="openDeleteProjectDialog(project)"
          >
            <v-icon>mdi-trash-can-outline</v-icon>
          </v-btn>
        </v-col>
      </v-row>

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
import UserSelect from "@/components/UserSelect"
import DeleteDialog from "@/components/DeleteDialog"

export default {
  name: "Projects",
  components: {
    CreateProjectDialog,
    EditProjectDialog,
    UserSelect,
    DeleteDialog,
  },

  data: () => ({
    createProjectDialog: false,

    editProjectDialog: false,
    projectToEdit: {},

    deleteProjectDialog: false,
    projectToDelete: {},

    projects: [],
    users: [],
    search: {},
  }),
  methods: {
    openDeleteProjectDialog(project) {
      this.projectToDelete = project
      this.deleteProjectDialog = true
    },
    deleteSelectedProject() {
      ProjectService.delete(this.projectToDelete)
          .then(() => this.refreshProjects())
    },
    clearFilters() {
      this.search = {}
      this.refreshProjects()
    },
    selectProject(project) {
      this.projectToEdit = {
        id: project.id,
        name: project.name,
        description: project.description,
        users: project.users,
      }
      this.editProjectDialog = true
    },
    getAllUserNames(project) {
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
  }
}
</script>