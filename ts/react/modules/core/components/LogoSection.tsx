import React from "react"

interface LogoSectionProps {
    lightMode?: boolean
    className?: string
    isHorizontal?: boolean
    logoSize?: number
}

const LogoSection: React.FC<LogoSectionProps> = ({
    lightMode = false,
    className = "",
    isHorizontal = false,
}: LogoSectionProps) => {
    return (
        <div
            className={`flex items-center
        ${isHorizontal ? "flex-row" : "flex-col"}
         gap-2 py-1`}
        >
            <img
                src="/brand/logo/raw.png"
                alt="Logo"
                className={`w-44 -my-12  ${lightMode ? "" : ""} ${className}`}
            />
        </div>
    )
}

export default LogoSection
