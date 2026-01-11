export const paths = {
  home: {
    path: "/admin/",
    getHref: () => "/admin/",
  },
  dashboard: {
    path: "/admin/dashboard",
    getHref: () => "/admin/dashboard",
  },
  addSketch: {
    path: "/admin/sketch/add",
    getHref: () => "/admin/sketch/add",
  },
  updateSketch: {
    path: "/admin/sketch/:id",
    getHref: (id: string | number) => `/admin/sketch/${id}`,
  },
  addPerson: {
    path: "/admin/person/add",
    getHref: () => "/admin/person/add",
  },
  addCharacter: {
    path: "/admin/character/add",
    getHref: () => "/admin/character/add",
  },
  sketches: {
    path: "/admin/sketches",
    getHref: () => "/admin/sketches",
  },
} as const;
