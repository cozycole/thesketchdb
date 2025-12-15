export const paths = {
  home: {
    path: "/",
    getHref: () => "/",
  },
  dashboard: {
    path: "/dashboard",
    getHref: () => "/dashboard",
  },
  addSketch: {
    path: "/sketch/add",
    getHref: () => "/sketch/add",
  },
  addPerson: {
    path: "/person/add",
    getHref: () => "/person/add",
  },
  addCharacter: {
    path: "/character/add",
    getHref: () => "/character/add",
  },
  sketches: {
    path: "/sketches",
    getHref: () => "/sketches",
  },
} as const;
