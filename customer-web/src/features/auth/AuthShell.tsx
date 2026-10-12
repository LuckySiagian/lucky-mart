import React from 'react';
import { Link } from 'react-router-dom';

export const authInput =
  'w-full rounded-xl border border-slate-700 bg-slate-800/80 py-2.5 pl-11 pr-4 text-sm text-white placeholder-slate-500 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500';

export const AuthShell: React.FC<{ title: string; subtitle: string; children: React.ReactNode }> = ({ title, subtitle, children }) => (
  <div className="flex min-h-screen items-center justify-center bg-slate-950 p-4 font-sans">
    <div className="w-full max-w-md space-y-6 rounded-3xl border border-slate-800 bg-slate-900 p-6 shadow-2xl sm:p-8">
      <div className="text-center">
        <Link to="/">
          <img src="/logo.png" alt="Lucky Mart" className="mx-auto h-14 w-auto object-contain" />
        </Link>
        <h2 className="mt-4 text-2xl font-extrabold text-white">{title}</h2>
        <p className="mt-1 text-sm text-slate-400">{subtitle}</p>
      </div>
      {children}
    </div>
  </div>
);
