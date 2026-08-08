import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

import TokenLogin from "./pages/TokenLogin";
import Welcome from "./pages/Welcome";
import SurveyWizard from "./pages/SurveyWizard";
import ThankYou from "./pages/ThankYou";
import AdminLogin from "./pages/admin/AdminLogin";
import AdminLayout from "./pages/admin/AdminLayout";
import AdminDashboard from "./pages/admin/AdminDashboard";
import ImportGraduates from "./pages/admin/ImportGraduates";
import SendInvites from './pages/admin/SendInvites';
import Settings from './pages/admin/Settings';

export default function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navigate to="/giris" replace />} />
                <Route path="/giris" element={<TokenLogin />} />
                <Route path="/hosgeldin" element={<Welcome />} />
                <Route path="/anket" element={<SurveyWizard />} />
                <Route path="/tesekkurler" element={<ThankYou />} />

                <Route path="/admin/giris" element={<AdminLogin />} />

                <Route path="/admin" element={<AdminLayout />}>
                    <Route index element={<AdminDashboard />} />
                    <Route path="mezun-ekle" element={<ImportGraduates />} />
                    <Route path="send-invites" element={<SendInvites />} />
                    <Route path="settings" element={<Settings />} />
                </Route>

                <Route path="*" element={<Navigate to="/giris" replace />} />
            </Routes>
        </BrowserRouter>
    );
}