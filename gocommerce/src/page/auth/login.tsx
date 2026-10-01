import {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useAuth} from '../../state/use-auth.ts';
import {ApiError} from '../../domain/common/api-client.tsx';

const LoginPage = () => {
    const [idToken, setIDToken] = useState('');
    const [error, setError] = useState('');
    const auth = useAuth();
    const navigate = useNavigate();

    return <main className='max-w-xl mx-auto bg-white rounded-lg shadow-sm p-6'>
        <h2 className='text-2xl font-semibold mb-2'>Sign in with Google</h2>
        <p className='text-sm text-gray-600 mb-4'>Pass the ID token returned by Google Identity Services. The backend verifies its signature, issuer, audience, and account status before creating a session.</p>
        <textarea value={idToken} onChange={event => setIDToken(event.target.value)} className='w-full border rounded p-2 min-h-28' placeholder='Google ID token'/>
        <button className='mt-4 bg-blue-600 text-white px-4 py-2 rounded' onClick={() => {
            setError('');
            auth.loginWithSSO(idToken).then(() => navigate('/checkout')).catch((reason: unknown) => {
                setError(reason instanceof ApiError ? `Login failed (${reason.status})` : 'Login failed');
            });
        }}>Continue</button>
        {error && <p className='text-red-600 mt-3'>{error}</p>}
    </main>;
};

export default LoginPage;
