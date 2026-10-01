import {Outlet} from 'react-router-dom';

const ProtectedRoute = () => {
    // const auth = useAuth();
    // const location = useLocation();
    return <Outlet/>;
};

export default ProtectedRoute;
