// Your hardcoded data
const incidents = [
    { id: 1, severity: "SEV1", title: "Payment", createdAt: "2026-06-02 17:20", resolvedAt: null },
    { id: 2, severity: "SEV3", title: "ABC", createdAt: "2026-06-01 17:20", resolvedAt: null },
    { id: 3, severity: "SEV2", title: "EFG", createdAt: "2026-06-02 10:30", resolvedAt: "2026-06-02 11:00" },
    { id: 4, severity: "SEV1", title: "XYZ", createdAt: "2026-06-03 08:00", resolvedAt: null },
    { id: 1, severity: "SEV1", title: "Payment", createdAt: "2026-06-01 17:20", resolvedAt: null },
];

const countBySeverity = (incidents) => {
    const results = {}
    incidents.forEach(inc => {
        results[inc.severity] = (results[inc.severity] || 0) + 1
    })
    return results
} 
const countOpen  = incidents.filter(inc => inc.resolvedAt === null).length
const countResolve = incidents.filter(inc => inc.resolvedAt !== null).length
const filteredOpen = incidents.filter(inc => inc.resolvedAt === null)
const filteredResolve = incidents.filter(inc => inc.resolvedAt !== null)

const sortBySeverityThenRecency = (incidents) => {
  const severityOrder = { SEV3: 3, SEV2: 2, SEV1: 1 };
  return [...incidents].sort((a, b) => {
    const sevDiff = severityOrder[b.severity] - severityOrder[a.severity];
    if (sevDiff !== 0) {
        return sevDiff;
    }
    return new Date(b.createdAt) - new Date(a.createdAt);
  });
};

console.log('countBySeverity("SEV1")', countBySeverity(incidents))
console.log('countOpen', countOpen)
console.log('countResolve', countResolve)
console.log(filteredOpen)
console.log(filteredResolve)

console.log(incidents)
console.log('countResolve', sortBySeverityThenRecency(incidents))
console.log(incidents)


function makeCounter() {
  let count = 0;
  return () => {
    count += 1;
    return count;
  };
}

const c = makeCounter();
console.log(c())
console.log(c())
console.log(c())