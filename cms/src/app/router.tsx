import { QueryClient, useQueryClient } from "@tanstack/react-query";
import { useMemo } from "react";
import { createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router/dom";

import { paths } from "@/config/paths";

import {
  default as AppRoot,
  ErrorBoundary as AppRootErrorBoundary,
} from "./routes/root";

const convert = (queryClient: QueryClient) => (m: any) => {
  const { clientLoader, clientAction, default: Component, ...rest } = m;
  return {
    ...rest,
    loader: clientLoader?.(queryClient),
    action: clientAction?.(queryClient),
    Component,
  };
};

export const createAppRouter = (queryClient: QueryClient) =>
  createBrowserRouter([
    {
      path: paths.home.path,
      element: <AppRoot />,
      ErrorBoundary: AppRootErrorBoundary,
      children: [
        {
          path: paths.home.path,
          lazy: () => import("./routes/dashboard").then(convert(queryClient)),
        },
        {
          path: paths.dashboard.path,
          lazy: () => import("./routes/dashboard").then(convert(queryClient)),
        },
        {
          path: paths.sketches.path,
          lazy: () => import("./routes/sketches").then(convert(queryClient)),
        },
        {
          path: paths.createSketch.path,
          lazy: () =>
            import("./routes/createSketch").then(convert(queryClient)),
        },
        {
          path: paths.updateSketch.path,
          lazy: () =>
            import("./routes/updateSketch").then(convert(queryClient)),
        },
        {
          path: paths.addPerson.path,
          lazy: () => import("./routes/addPerson").then(convert(queryClient)),
        },
        {
          path: paths.addCharacter.path,
          lazy: () =>
            import("./routes/addCharacter").then(convert(queryClient)),
        },
      ],
    },
    {
      path: "*",
      lazy: () => import("./routes/not-found").then(convert(queryClient)),
    },
  ]);

export const AppRouter = () => {
  const queryClient = useQueryClient();

  const router = useMemo(() => createAppRouter(queryClient), [queryClient]);

  return <RouterProvider router={router} />;
};
