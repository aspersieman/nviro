document.addEventListener('alpine:init', () => {
  Alpine.store('main', {
    activeTab: 'environments',
    showProjectAddModal: false,
    projects: [],
    projectId: 'all',
    projectName: null,
    projectSearch: '',
    projectModalId: null,
    showEnvironmentAddModal: false,
    environments: [],
    showEnironmentHistory: null,
    environmentHistory: [],
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
      this.projects = [];
      const response = await fetch('/api/projects');
      const projects = await response.json();
      this.projects = projects;
      this.projectList();
    },

    createProject(name) {
      if (name) {
        fetch('/api/projects', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            'name': name
          })
        }).then(() => {
          this.clearProjectAddModal();
          this.showProjectAddModal = false;
          this.getProjects();
        })
      }
      return
    },

    updateProject(id, name) {
      if (id && name) {
        fetch('/api/projects/' + id, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            'name': name
          })
        }).then(() => {
          this.clearProjectAddModal();
          this.showProjectAddModal = false;
          this.getProjects();
        });
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
          this.getProjects();
        });
      }
    },

    showProjectAdd() {
      this.clearProjectAddModal();
      this.showProjectAddModal = true;
    },

    showProjectEdit(id, name) {
      this.clearProjectAddModal();
      this.projectModalId = id;
      this.projectName = name;
      this.showProjectAddModal = true;
    },

    clearProjectSearch() {
      this.projectSearch = '';
    },

    clearProjectAddModal() {
      this.projectModalId = null;
      this.projectName = null;
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

    async updateEnvironment(
      id,
      name,
      content,
      project_id
    ) {
      if (name && content && project_id) {
        const response = await fetch('/api/environments/' + id, {
          method: 'PUT',
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
      this.environments = [];
      const paramDeleted = 'false';
      const url = (
        '/api/environments?' +
        new URLSearchParams({ deleted: paramDeleted }).toString()
      );
      const response = await fetch(url);
      const environments = await response.json();
      this.environments = environments;
      this.environmentList();
    },

    async getEnvironmentHistory(id, name, project_id) {
      this.environmentHistory = [];
      const url = (
        '/api/environments?' +
        new URLSearchParams({
          id,
          name,
          project_id,
          deleted: true
        }).toString()
      );
      const response = await fetch(url);
      const environments = await response.json();
      this.environmentHistory = environments;
      this.showEnironmentHistory = id;
    },

    async showEnironmentHistoryList(id, name, project_id) {
      if (this.showEnironmentHistory != id) {
        this.getEnvironmentHistory(id, name, project_id);
      } else {
        this.showEnironmentHistory = null;
      }
      
    },

    deleteEnvironment(id, force) {
      params = force ? '?force=true' : '?force=false';
      if (id) {
        fetch('/api/environments/' + id + params, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          }
        }).then(() => {
          this.getEnvironments();
          this.environmentHistory = this.environmentHistory.filter(environment => {
            return environment.id !== id; 
          })
          if (this.showEnironmentHistory.length === 0) {
            this.showEnironmentHistory = null;
          }
        });
      }
    },

    clearEnvironmentSearch() {
      this.environmentSearch = '';
      this.getEnvironments();
    },

    showEnvironmentAdd() {
      this.clearEnvironmentAddModal();
      this.showEnvironmentAddModal = true;
    },

    hideEnvironmentAdd() {
      this.clearEnvironmentAddModal();
      this.showEnvironmentAddModal = false;
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
      this.clearEnvironmentAddModal();
      this.environmentId = id;
      this.environmentName = name;
      this.environmentProjectId = projectId;
      this.environmentContent = content;
      this.showEnvironmentAddModal = true;
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
      switch (tabName) {
        case 'environments':
          this.getEnvironments();
          break;
        case 'projects':
          this.getProjects();
          break;
        default:
          break;
      }
      this.activeTab = tabName;
    },
  })
})
