import { AppProvider } from "@/app/provider";
import { AppRouter } from "@/app/router";

const App = () => {
  return (
    <AppProvider>
      <AppRouter />
    </AppProvider>
  );
};

export default App;
