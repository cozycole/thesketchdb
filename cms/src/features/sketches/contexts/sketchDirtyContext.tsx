import {
  createContext,
  useContext,
  useMemo,
  useCallback,
  useState,
} from "react";

type DirtyContextValue = {
  dirtyMap: Record<string, boolean>;
  setDirty: (key: string, dirty: boolean) => void;
  hasUnsavedChanges: boolean;
};

const DirtyContext = createContext<DirtyContextValue | null>(null);

export function SketchDirtyProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [dirtyMap, setDirtyMap] = useState<Record<string, boolean>>({});
  const setDirty = useCallback((key: string, dirty: boolean) => {
    setDirtyMap((prev) => {
      if (prev[key] === dirty) return prev;

      return {
        ...prev,
        [key]: dirty,
      };
    });
  }, []);

  const value = useMemo(() => {
    return {
      dirtyMap,
      setDirty,
      hasUnsavedChanges: Object.values(dirtyMap).some(Boolean),
    };
  }, [dirtyMap, setDirty]);

  return (
    <DirtyContext.Provider value={value}>{children}</DirtyContext.Provider>
  );
}

export function useSketchDirty() {
  const ctx = useContext(DirtyContext);
  if (!ctx)
    throw new Error("useSketchDirty must be used inside SketchDirtyProvider");
  return ctx;
}
