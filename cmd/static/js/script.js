document.addEventListener('alpine:init', () => {
  Alpine.store('main', {
    showEnvironmentAddModal: false,
    environments: [],
    environmentSearch: null,
    projects: [],
    projectId: null,
    projectName: null,
    environmentName: null,
    environmentProjectId: null,
    environmentContent: null,

    init() {
      this.getProjects();
      this.getEnvironments();
    },

    async getEnvironments() {
      const response = await fetch("/api/environments");
      const environments = await response.json();
      this.environments = environments;
    },

    async getProjects() {
      const response = await fetch("/api/projects");
      const projects = await response.json();
      this.projects = projects;
    },

    async createProject(name) {
      if (name) {
        const response = await fetch("/api/projects", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({
            "name": name
          })
        });
        this.clearProjectAddModal();
        this.showProjectAddModal = false;
        this.getProjects();
        return response
      }
      return
    },

    async createEnvironment(
      name,
      content,
      project_id
    ) {
      if (name && content && project_id) {
        const response = await fetch("/api/environments", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({
            "name": name,
            "content": content,
            "project_id": parseInt(project_id)
          })
        });
        this.clearEnvironmentAddModal();
        this.showEnvironmentAddModal = false;
        this.getEnvironments();
        return response
      }
      return
    },

    deleteEnvironment(id) {
      if (id) {
        fetch("/api/environments/" + id, {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json"
          }
        }).then(() => {
          window.location.reload();
        });
      }
    },

    async downloadContent(name, content) {
      const link = document.createElement("a");
      link.download = name.toString();

      const file = new Blob([content], { type: 'text/plain;charset=utf-8' });

      link.href = URL.createObjectURL(file);
      link.click();
      URL.revokeObjectURL(link.href);
    },

    filterEnvironmentsName() {
      this.environments = this.environments.filter(environment => {
        return environment.name.toLowerCase().includes(this.environmentSearch.toLowerCase()); 
      })
    },

    filterEnvironmentsProject() {
      console.log('filterEnvironmentsProject: ', this.projectId);
      if (this.projectId === "all") {
        this.getEnvironments();
        return
      }
      const self = this;
      this.getEnvironments().then(() => {
        self.environments = self.environments.filter((environment) => {
          const include = parseInt(environment.project_id) === parseInt(self.projectId)
          console.log('filterEnvironmentsProject: ', environment.project_name, environment.name, 'env project id ', environment.project_id, include);
          return include; 
        })
      })
    },

    clearSearch() {
      this.environmentSearch = null;
      this.getEnvironments();
    },

    clearEnvironmentAddModal() {
      this.environmentName = null;
      this.environmentProjectId = null;
      this.environmentContent = null;
    },

    clearProjectAddModal() {
      this.projectName = null;
    }
  })
})
