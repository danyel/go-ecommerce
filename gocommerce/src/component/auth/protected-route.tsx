import {Navigate, Outlet, useLocation} from 'react-router-dom';
import {useAuth} from '../../state/use-auth.ts';

const ProtectedRoute = () => {
    const auth = useAuth();
    const location = useLocation();
    return auth.isAuthenticated ? <Outlet/> : <Navigate to='/login' replace state={{from: location.pathname}}/>;
};

export default ProtectedRoute;
