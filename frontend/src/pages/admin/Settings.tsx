import React, { useState, useEffect } from 'react';

export default function Settings() {
    const [settings, setSettings] = useState({
        smtp_host: '',
        smtp_port: '',
        smtp_user: '',
        smtp_pass: '',
        sms_api_key: '',
        sms_api_url: ''
    });

    // Projende admin token'ýný nerede tutuyorsan (localStorage vb.) ona göre burayý güncelle
    const getAuthToken = () => localStorage.getItem('token');

    // Sayfa yüklendiðinde mevcut ayarlarý çek
    useEffect(() => {
        const token = getAuthToken();

        fetch('http://localhost:8080/api/v1/admin/settings', {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        })
            .then(res => res.json())
            .then(response => {
                // Backend middleware'inden dönen success envelope yapýsýna göre (response.data)
                if (response.data) {
                    setSettings(prev => ({ ...prev, ...response.data }));
                }
            })
            .catch(err => console.error("Ayarlar yüklenirken hata oluþtu:", err));
    }, []);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSettings({ ...settings, [e.target.name]: e.target.value });
    };

    const handleSave = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const token = getAuthToken();

        try {
            const res = await fetch('http://localhost:8080/api/v1/admin/settings', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(settings)
            });

            if (res.ok) {
                alert('Ayarlar baþarýyla güncellendi!');
            } else {
                alert('Kaydedilirken bir hata oluþtu.');
            }
        } catch (error) {
            console.error("Kaydetme hatasý:", error);
            alert('Sunucuya ulaþýlamadý.');
        }
    };

    return (
        <div className="p-6 bg-gray-900 text-white rounded-lg max-w-3xl mx-auto mt-10 shadow-lg">
            <h2 className="text-2xl font-bold mb-6 text-orange-500 border-b border-gray-700 pb-2">
                Sistem Ayarlarý (SMTP & SMS)
            </h2>
            <form onSubmit={handleSave} className="space-y-6">

                {/* SMTP Ayarlarý Kartý */}
                <div className="p-5 bg-gray-800 border border-gray-700 rounded-md shadow-sm">
                    <h3 className="text-lg mb-4 font-semibold text-gray-200">E-posta (SMTP) Ayarlarý</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="flex flex-col">
                            <label className="text-sm text-gray-400 mb-1">SMTP Sunucusu</label>
                            <input className="p-2 bg-gray-700 rounded border border-gray-600 focus:outline-none focus:border-orange-500" name="smtp_host" placeholder="Örn: smtp.mersin.edu.tr" value={settings.smtp_host || ''} onChange={handleChange} />
                        </div>
                        <div className="flex flex-col">
                            <label className="text-sm text-gray-400 mb-1">Port</label>
                            <input className="p-2 bg-gray-700 rounded border border-gray-600 focus:outline-none focus:border-orange-500" name="smtp_port" placeholder="Örn: 587" value={settings.smtp_port || ''} onChange={handleChange} />
                        </div>
                        <div className="flex flex-col">
                            <label className="text-sm text-gray-400 mb-1">E-posta Adresi (Kullanýcý)</label>
                            <input className="p-2 bg-gray-700 rounded border border-gray-600 focus:outline-none focus:border-orange-500" name="smtp_user" placeholder="E-posta" value={settings.smtp_user || ''} onChange={handleChange} />
                        </div>
                        <div className="flex flex-col">
                            <label className="text-sm text-gray-400 mb-1">Þifre</label>
                            <input className="p-2 bg-gray-700 rounded border border-gray-600 focus:outline-none focus:border-orange-500" name="smtp_pass" type="password" placeholder="Þifre" value={settings.smtp_pass || ''} onChange={handleChange} />
                        </div>
                    </div>
                </div>

                {/* SMS Ayarlarý Kartý */}
                <div className="p-5 bg-gray-800 border border-gray-700 rounded-md shadow-sm">
                    <h3 className="text-lg mb-4 font-semibold text-gray-200">SMS API Ayarlarý</h3>
                    <div className="flex flex-col space-y-4">
                        <div className="flex flex-col">
                            <label className="text-sm text-gray-400 mb-1">API URL</label>
                            <input className="p-2 bg-gray-700 rounded border border-gray-600 focus:outline-none focus:border-orange-500" name="sms_api_url" placeholder="https://api.sms-saglayici.com/send" value={settings.sms_api_url || ''} onChange={handleChange} />
                        </div>
                        <div className="flex flex-col">
                            <label className="text-sm text-gray-400 mb-1">API Anahtarý</label>
                            <input className="p-2 bg-gray-700 rounded border border-gray-600 focus:outline-none focus:border-orange-500" name="sms_api_key" type="password" placeholder="Gizli Anahtar" value={settings.sms_api_key || ''} onChange={handleChange} />
                        </div>
                    </div>
                </div>

                <button type="submit" className="w-full bg-orange-600 hover:bg-orange-500 text-white font-bold py-3 px-4 rounded transition-colors duration-200 flex justify-center items-center">
                    Ayarlarý Kaydet
                </button>
            </form>
        </div>
    );
}