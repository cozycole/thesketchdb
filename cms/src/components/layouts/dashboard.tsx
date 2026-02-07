import { Home, Database, Users } from "lucide-react";
import { NavLink, useNavigate, useNavigation } from "react-router";

import { paths } from "@/config/paths";

import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/ui/appSidebar";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const addMenuItems = [
  {
    title: "Sketch",
    path: paths.createSketch.getHref(),
  },
  {
    title: "Person",
    path: paths.addPerson.getHref(),
  },
  {
    title: "Character",
    path: paths.addCharacter.getHref(),
  },
  {
    title: "Creator",
    path: paths.createSketch.getHref(),
  },
  {
    title: "Show",
    path: paths.createSketch.getHref(),
  },
  {
    title: "Category / Tag",
    path: paths.createSketch.getHref(),
  },
];

function AddItemDropdown() {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button>Create +</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-44" align="end">
        {addMenuItems.map((item) => (
          <NavLink to={item.path} key={item.title}>
            <DropdownMenuItem className="flex flex-col">
              {item.title}
            </DropdownMenuItem>
          </NavLink>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="">
      <SidebarProvider>
        <AppSidebar />
        <main className="flex-1 items-start">
          <header className="sticky top-0 z-30 flex h-14 items-center justify-between gap-4 border-b bg-background p-4">
            <SidebarTrigger />
            <AddItemDropdown />
          </header>
          {children}
        </main>
      </SidebarProvider>
    </div>
  );
}
