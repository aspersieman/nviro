const filterEnvironments = () => {
  const trs = document.querySelectorAll('#tbl-environments tbody tr')
  console.log(trs)
  const filter = document.querySelector('#txt-environment-search').value
  console.log('filterEnvironments: ', filter)
  const regex = new RegExp(filter, 'i')
  const isFoundInTds = td => regex.test(td.innerHTML)
  const isFound = childrenArr => childrenArr.some(isFoundInTds)
  const setTrStyleDisplay = ({ style, children }) => {
    style.display = isFound([
      ...children // <-- All columns
    ]) ? '' : 'none' 
  }

  trs.forEach(setTrStyleDisplay)
}

const getEnvironments = async () => {
  const response = await fetch("/api/environments");
  const environments = await response.json();
  return environments;
}

const getProjects = async () => {
  const response = await fetch("/api/projects");
  const projects = await response.json();
  return projects;
}

const createEnvironment = async (
  name,
  content,
  project_id
) => {
  const response = await fetch("/api/environments", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      name,
      content,
      project_id
    })
  });
  const environment = await response.json();
  return environment; 
}

function alpineInstance() {
  return {
    environments: [],
    projects: [],
    showEnvironmentAddModal: false,
    getEnvironments,
    getProjects,
    createEnvironment
  }    
}
