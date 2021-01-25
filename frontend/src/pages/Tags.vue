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
          hide-default-footer
      >
        <template v-slot:item.color="{ item }">
          <v-icon v-bind:color="item.color">mdi-format-color-highlight</v-icon>
        </template>

        <template v-slot:item.actions="{ item }">
          <div class="d-flex justify-end">
            <v-btn
                icon
                @click="editTag(item)"
            >
              <v-icon>mdi-pencil-outline</v-icon>
            </v-btn>
            <v-btn
                icon
                @click="openDeleteTagDialog(item)"
            >
              <v-icon>mdi-trash-can-outline</v-icon>
            </v-btn>
          </div>
        </template>

      </v-data-table>

      <create-tag-dialog v-model="showCreateTagDialog"
                         @refreshTags="refreshTags"/>

      <edit-tag-dialog v-model="showEditTagDialog"
                       :tag="tagToEdit"
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

export default {
  name: "Tags",
  components: {
    CreateTagDialog,
    EditTagDialog,
    DeleteDialog,
  },

  data: () => ({
    showCreateTagDialog: false,
    showEditTagDialog: false,
    showDeleteTagDialog: false,

    tagToEdit: {},
    tagToDelete: {},

    tagsHeaders: [
      {text: 'NAME', value: 'name'},
      {text: 'DESCRIPTION', value: 'description'},
      {text: 'COLOR', value: 'color', sortable: false},
      {text: '', value: 'actions', sortable: false},
    ],
    tags: [],
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
    deleteTag() {
      TagService.delete(this.tagToDelete)
          .then(() => this.refreshTags())
    },
    refreshTags() {
      TagService.findAll()
          .then(response => this.tags = response.data)
    },
  },
  created() {
    this.refreshTags()
  }
}
</script>