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
