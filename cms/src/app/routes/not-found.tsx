import { Link } from "react-router";
import { Button } from "@/components/ui/button";
import { paths } from "@/config/paths";

const NotFoundRoute = () => {
  return (
    <div className="mt-52 flex flex-col items-center gap-2 font-semibold">
      <h1>404 - Not Found</h1>
      <p className="mb-6">
        Sorry, the page you are looking for does not exist.
      </p>
      <Link to={paths.home.getHref()} replace>
        <Button>Go to Home</Button>
      </Link>
    </div>
  );
};

export default NotFoundRoute;
