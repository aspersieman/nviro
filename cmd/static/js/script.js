document.addEventListener('alpine:init', () => {
  Alpine.store('main', {
    activeTab: 'environments',
    showProjectAddModal: false,
    projects: [],
    projectId: 'all',
    projectName: null,
    projectSearch: '',
    showEnvironmentAddModal: false,
    environments: [],
    environmentListDeleted: false,
    environmentSearch: '',
    environmentId: null,
    environmentName: null,
    environmentProjectId: null,
    environmentContent: null,

    init() {
      this.getProjects();
      this.getEnvironments();
    },

    projectList() {
      return this.projects.filter(project => {
        return this.projectSearch === null || project.name.toLowerCase().includes(this.projectSearch.toLowerCase()); 
      })
    },


    async getProjects() {
      const response = await fetch('/api/projects');
      const projects = await response.json();
      this.projects = projects;
    },

    async createProject(name) {
      if (name) {
        const response = await fetch('/api/projects', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            'name': name
          })
        });
        this.clearProjectAddModal();
        this.showProjectAddModal = false;
        this.getProjects();
        return response
      }
      return
    },

    deleteProject(id) {
      if (id) {
        fetch('/api/projects/' + id, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          }
        }).then(() => {
          window.location.reload();
        });
      }
    },

    async createEnvironment(
      name,
      content,
      project_id
    ) {
      if (name && content && project_id) {
        const response = await fetch('/api/environments', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            'name': name,
            'content': content,
            'project_id': parseInt(project_id)
          })
        });
        this.clearEnvironmentAddModal();
        this.showEnvironmentAddModal = false;
        this.getEnvironments();
        return response
      } else {
        console.log('createEnvironment: stopped: ', name, content, project_id);
      }
      return
    },

    environmentList() {
      return this.environments.filter(environment => {
        return environment.name.toLowerCase().includes(this.environmentSearch.toLowerCase()) && (
          this.projectId === 'all' || parseInt(environment.project_id) === parseInt(this.projectId)
        );
      })
    },

    async getEnvironments() {
      const paramDeleted = this.environmentListDeleted ? 'true' : 'false';
      const url = (
        '/api/environments?' +
        new URLSearchParams({ deleted: paramDeleted }).toString()
      );
      const response = await fetch(url);
      const environments = await response.json();
      this.environments = environments;
    },

    deleteEnvironment(id) {
      if (id) {
        fetch('/api/environments/' + id, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          }
        }).then(() => {
          window.location.reload();
        });
      }
    },

    async downloadContent(name, content) {
      const link = document.createElement('a');
      link.download = name.toString();

      const file = new Blob([content], { type: 'text/plain;charset=utf-8' });

      link.href = URL.createObjectURL(file);
      link.click();
      URL.revokeObjectURL(link.href);
    },

    setActiveTab(tabName) {
      this.activeTab = tabName;
    },

    clearProjectSearch() {
      this.projectSearch = '';
    },

    clearProjectAddModal() {
      this.projectName = null;
    },

    clearEnvironmentSearch() {
      this.environmentSearch = '';
      this.getEnvironments();
    },

    clearEnvironmentAddModal() {
      this.environmentId = null;
      this.environmentName = null;
      this.environmentProjectId = null;
      this.environmentContent = null;
    },
    
    environmentAddModalEdit(
      id,
      name,
      projectId,
      content
    ) {
      this.environmentId = id;
      this.environmentName = name;
      this.environmentProjectId = projectId;
      this.environmentContent = content;
      console.log('environmentAddModalEdit: ', id, name, projectId, content);
      this.showEnvironmentAddModal = true;
    }
  })
})
