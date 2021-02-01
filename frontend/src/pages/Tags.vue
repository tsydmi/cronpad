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
                @click="showCreateTagDialog = true"
                v-bind="attrs"
                v-on="on"
            >
              <v-icon>mdi-plus-circle</v-icon>
            </v-btn>
          </template>
          <span>ADD TAG</span>
        </v-tooltip>
      </v-col>

      <v-data-table
          :headers="tagsHeaders"
          :items="tags"
          group-by="project"
          hide-default-footer
      >
        <template v-slot:item.parent="{ item }">
          <span> {{ getNameById(tags, item.parent) }} </span>
        </template>
        <template v-slot:item.project="{ item }">
          <span> {{ getNameById(projects, item.project) }} </span>
        </template>
        <template v-slot:item.color="{ item }">
          <v-icon :color="item.color">mdi-format-color-highlight</v-icon>
        </template>

        <template v-slot:item.actions="{ item }">
          <div class="d-flex justify-end">
            <v-btn
                icon
                :disabled="item.basic && !hasAdminRole"
                @click="editTag(item)"
            >
              <v-icon>mdi-pencil-outline</v-icon>
            </v-btn>
            <v-btn
                icon
                :disabled="item.basic && !hasAdminRole"
                @click="openDeleteTagDialog(item)"
            >
              <v-icon>mdi-trash-can-outline</v-icon>
            </v-btn>
          </div>
        </template>

        <template v-slot:group.header="{ group, toggle }">
          <td class="pa-2">
            <v-btn
                @click="toggle"
                icon
            >
              <v-icon>fa-minus</v-icon>
            </v-btn>
            <span class="pl-3">
              {{ group ? `Project: ${getNameById(projects, group)}` : 'Basic tags' }}
            </span>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
        </template>

      </v-data-table>

      <create-tag-dialog v-model="showCreateTagDialog"
                         :tags="basicTags"
                         :projects="projects"
                         :hasAdminRole="hasAdminRole"
                         @refreshTags="refreshTags"/>

      <edit-tag-dialog v-model="showEditTagDialog"
                       :tag="tagToEdit"
                       :tags="basicTags"
                       :projects="projects"
                       :hasAdminRole="hasAdminRole"
                       @refreshTags="refreshTags"/>

      <delete-dialog v-model="showDeleteTagDialog"
                     @delete="deleteTag"
                     :item-name="tagToDelete.name"/>
    </v-card>
  </v-container>
</template>

<script>
import TagService from "@/service/TagService"
import CreateTagDialog from "@/components/tags/CreateTagDialog"
import EditTagDialog from "@/components/tags/EditTagDialog"
import DeleteDialog from "@/components/DeleteDialog"
import cloneDeep from 'clone-deep'
import ProjectService from "@/service/ProjectService"

export default {
  name: "Tags",
  components: {
    CreateTagDialog,
    EditTagDialog,
    DeleteDialog,
  },

  data: () => ({
    hasAdminRole: false,

    showCreateTagDialog: false,
    showEditTagDialog: false,
    showDeleteTagDialog: false,

    tagToEdit: {},
    tagToDelete: {},

    tagsHeaders: [
      {text: 'NAME', value: 'name'},
      {text: 'DESCRIPTION', value: 'description'},
      {text: 'PARENT', value: 'parent'},
      {text: 'PROJECT', value: 'project'},
      {text: 'COLOR', value: 'color', sortable: false},
      {text: '', value: 'actions', sortable: false},
    ],
    tags: [],
    basicTags: [],
    projects: [],
  }),
  methods: {
    editTag(tag) {
      this.tagToEdit = cloneDeep(tag)
      this.showEditTagDialog = true
    },
    openDeleteTagDialog(tag) {
      this.tagToDelete = tag
      this.showDeleteTagDialog = true
    },
    getNameById(array, id) {
      let element = array.find(e => e.id === id);
      if (element) {
        return element.name
      }
      return id
    },
    deleteTag() {
      TagService.delete(this.tagToDelete)
          .then(() => this.refreshTags())
    },
    refreshTags() {
      TagService.findAll()
          .then(response => {
            this.tags = response.data
            this.basicTags = response.data.filter(t => t.basic)
          })
    },
    refreshProjects() {
      ProjectService.findAll(this.hasAdminRole)
          .then(response => this.projects = response.data)
    },
  },
  created() {
    this.hasAdminRole = this.$keycloak.hasRealmRole('admin')
    this.refreshTags()
    this.refreshProjects()
  },
}
</script>